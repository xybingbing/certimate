﻿package tencentcloudcos

import (
	"context"
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tcSsl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/usual2970/certimate/internal/pkg/core/deployer"
	"github.com/usual2970/certimate/internal/pkg/core/logger"
	"github.com/usual2970/certimate/internal/pkg/core/uploader"
	uploadersp "github.com/usual2970/certimate/internal/pkg/core/uploader/providers/tencentcloud-ssl"
)

type DeployerConfig struct {
	// 腾讯云 SecretId。
	SecretId string `json:"secretId"`
	// 腾讯云 SecretKey。
	SecretKey string `json:"secretKey"`
	// 腾讯云地域。
	Region string `json:"region"`
	// 存储桶名。
	Bucket string `json:"bucket"`
	// 自定义域名（不支持泛域名）。
	Domain string `json:"domain"`
}

type DeployerProvider struct {
	config      *DeployerConfig
	logger      logger.Logger
	sdkClient   *tcSsl.Client
	sslUploader uploader.Uploader
}

var _ deployer.Deployer = (*DeployerProvider)(nil)

func NewDeployer(config *DeployerConfig) (*DeployerProvider, error) {
	if config == nil {
		panic("config is nil")
	}

	client, err := createSdkClient(config.SecretId, config.SecretKey, config.Region)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create sdk clients")
	}

	uploader, err := uploadersp.NewUploader(&uploadersp.UploaderConfig{
		SecretId:  config.SecretId,
		SecretKey: config.SecretKey,
	})
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to create ssl uploader")
	}

	return &DeployerProvider{
		config:      config,
		logger:      logger.NewNilLogger(),
		sdkClient:   client,
		sslUploader: uploader,
	}, nil
}

func (d *DeployerProvider) WithLogger(logger logger.Logger) *DeployerProvider {
	d.logger = logger
	return d
}

func (d *DeployerProvider) Deploy(ctx context.Context, certPem string, privkeyPem string) (*deployer.DeployResult, error) {
	if d.config.Bucket == "" {
		return nil, errors.New("config `bucket` is required")
	}
	if d.config.Domain == "" {
		return nil, errors.New("config `domain` is required")
	}

	// 上传证书到 SSL
	upres, err := d.sslUploader.Upload(ctx, certPem, privkeyPem)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to upload certificate file")
	}

	d.logger.Logt("certificate file uploaded", upres)

	// 证书部署到 COS 实例
	// REF: https://cloud.tencent.com/document/product/400/91667
	deployCertificateInstanceReq := tcSsl.NewDeployCertificateInstanceRequest()
	deployCertificateInstanceReq.CertificateId = common.StringPtr(upres.CertId)
	deployCertificateInstanceReq.ResourceType = common.StringPtr("cos")
	deployCertificateInstanceReq.Status = common.Int64Ptr(1)
	deployCertificateInstanceReq.InstanceIdList = common.StringPtrs([]string{fmt.Sprintf("%s#%s#%s", d.config.Region, d.config.Bucket, d.config.Domain)})
	deployCertificateInstanceResp, err := d.sdkClient.DeployCertificateInstance(deployCertificateInstanceReq)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to execute sdk request 'ssl.DeployCertificateInstance'")
	}

	d.logger.Logt("已部署证书到云资源实例", deployCertificateInstanceResp.Response)

	return &deployer.DeployResult{}, nil
}

func createSdkClient(secretId, secretKey, region string) (*tcSsl.Client, error) {
	credential := common.NewCredential(secretId, secretKey)
	client, err := tcSsl.NewClient(credential, region, profile.NewClientProfile())
	if err != nil {
		return nil, err
	}

	return client, nil
}
