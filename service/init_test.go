package service

import (
	"testing"
	"time"

	specV1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	v1 "github.com/baetyl/baetyl-go/v2/spec/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/baetyl/baetyl-cloud/v2/common"
	"github.com/baetyl/baetyl-cloud/v2/mock/service"
	"github.com/baetyl/baetyl-cloud/v2/models"
)

func TestInitService_GetResource(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	tp := service.NewMockTemplateService(mockCtl)
	ns := service.NewMockNodeService(mockCtl)
	aus := service.NewMockAuthService(mockCtl)
	as := InitServiceImpl{}
	as.TemplateService = tp
	as.NodeService = ns
	as.AuthService = aus

	// good case : metrics
	tp.EXPECT().GetTemplate(templateKubeAPIMetricsYaml).Return("metrics", nil).Times(1)
	res, _ := as.GetResource(templateKubeAPIMetricsYaml, "", "", nil)
	assert.Equal(t, res, []byte("metrics"))

	// good case : local_path_storage
	tp.EXPECT().GetTemplate(templateKubeLocalPathStorageYaml).Return("local-path-storage", nil).Times(1)
	res, _ = as.GetResource(templateKubeLocalPathStorageYaml, "", "", nil)
	assert.Equal(t, res, []byte("local-path-storage"))

	// good case : setup
	tp.EXPECT().ParseTemplate(templateInitSetupShell, gomock.Any()).Return([]byte("shell"), nil).Times(1)
	res, _ = as.GetResource(templateInitSetupShell, "", "", nil)
	assert.Equal(t, res, []byte("shell"))

	// bad case : not found
	_, err := as.GetResource("dummy", "", "", nil)
	assert.EqualError(t, err, "The (resource) resource (dummy) is not found.")
}

func TestInitService_getInitYaml(t *testing.T) {
	info := map[string]interface{}{
		InfoKind:      "123",
		InfoName:      "n0",
		InfoNamespace: "default",
		InfoTimestamp: time.Now().Unix(),
		InfoExpiry:    60 * 60 * 24 * 3650,
	}
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	auth := service.NewMockAuthService(mockCtl)
	as.AuthService = auth
	res, err := as.getInitYaml(info, "kube")
	assert.Error(t, err, common.ErrRequestParamInvalid)
	assert.Nil(t, res)
}

func TestInitService_getSyncCert(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	secret := service.NewMockSecretService(mockCtl)
	as.SecretService = secret

	app1 := &specV1.Application{
		Name:      "baetyl-core",
		Namespace: "default",
	}

	secret.EXPECT().Get("default", "", "").Return(nil, nil).Times(1)
	res, err := as.getNodeCert(app1)
	assert.Error(t, err, common.ErrResourceNotFound)
	assert.Nil(t, res)
}

func TestInitService_GenCmd(t *testing.T) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	sAuth := service.NewMockAuthService(mockCtl)
	sTemplate := service.NewMockTemplateService(mockCtl)
	as := InitServiceImpl{}
	as.AuthService = sAuth
	as.TemplateService = sTemplate

	info := map[string]interface{}{
		InfoKind:      "kind",
		InfoName:      "name",
		InfoNamespace: "ns",
		InfoExpiry:    CmdExpirationInSeconds,
		InfoTimestamp: time.Now().Unix(),
	}
	expect := "curl -skfL 'https://1.2.3.4:9003/v1/active/setup.sh?token=tokenexpect' -osetup.sh && sh setup.sh"

	sAuth.EXPECT().GenToken(info).Return("tokenexpect", nil).Times(1)
	sTemplate.EXPECT().Execute("setup-command", templateCommand, gomock.Any()).Return([]byte(expect), nil).Times(1)

	res, err := as.GenCmd("kind", "ns", "name")
	assert.NoError(t, err)
	assert.Equal(t, res, expect)
}

func TestInitService_getDesireAppInfo(t *testing.T) {
	as := InitServiceImpl{}
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	node := service.NewMockNodeService(mockCtl)
	app := service.NewMockApplicationService(mockCtl)
	as.NodeService = node
	as.AppCombinedService = &AppCombinedService{
		App: app,
	}

	Desire := &specV1.Desire{
		"sysapps": []specV1.AppInfo{{
			Name:    "baetyl-core-node01",
			Version: "123",
		}},
	}
	app1 := &specV1.Application{
		Name:      "baetyl-core",
		Namespace: "default",
	}
	node.EXPECT().GetDesire("default", "node01").Return(Desire, nil).Times(1)
	app.EXPECT().Get("default", "baetyl-core-node01", "").Return(app1, nil).Times(1)

	res, err := as.getCoreAppFromDesire("default", "node01")
	assert.NoError(t, err)
	assert.Equal(t, res, app1)
}

func TestInitService_GenApps(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	sApp := service.NewMockApplicationService(mock)
	sConfig := service.NewMockConfigService(mock)
	sSecret := service.NewMockSecretService(mock)
	sTemplate := service.NewMockTemplateService(mock)
	sNode := service.NewMockNodeService(mock)
	sAuth := service.NewMockAuthService(mock)
	sPKI := service.NewMockPKIService(mock)

	is := InitServiceImpl{}
	is.TemplateService = sTemplate
	is.NodeService = sNode
	is.AuthService = sAuth
	is.PKI = sPKI
	is.AppCombinedService = &AppCombinedService{
		App:    sApp,
		Config: sConfig,
		Secret: sSecret,
	}

	cert := &models.PEMCredential{
		CertPEM: []byte("CertPEM"),
		KeyPEM:  []byte("KeyPEM"),
		CertId:  "CertId",
	}

	config := &v1.Configuration{
		Namespace: "ns",
		Name:      "config",
	}

	secret := &v1.Secret{
		Namespace: "ns",
		Name:      "secret",
	}

	app := &v1.Application{
		Namespace: "ns",
		Name:      "app",
	}

	sTemplate.EXPECT().UnmarshalTemplate("baetyl-core-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-core-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-function-conf.yml", gomock.Any(), gomock.Any()).Return(nil)
	sTemplate.EXPECT().UnmarshalTemplate("baetyl-function-app.yml", gomock.Any(), gomock.Any()).Return(nil)
	sPKI.EXPECT().SignClientCertificate("ns.abc", gomock.Any()).Return(cert, nil)
	sPKI.EXPECT().GetCA().Return([]byte("RootCA"), nil)
	sConfig.EXPECT().Create("ns", gomock.Any()).Return(config, nil).Times(2)
	sSecret.EXPECT().Create("ns", gomock.Any()).Return(secret, nil).Times(1)
	sApp.EXPECT().Create("ns", gomock.Any()).Return(app, nil).Times(2)

	out, err := is.GenApps("ns", "abc", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(out))
}