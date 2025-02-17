package applicant

import (
	"fmt"

	"github.com/go-acme/lego/v4/challenge"

	"github.com/usual2970/certimate/internal/domain"
	pACMEHttpReq "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/acmehttpreq"
	pAliyun "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/aliyun"
	pAWSRoute53 "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/aws-route53"
	pAzureDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/azure-dns"
	pCloudflare "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/cloudflare"
	pClouDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/cloudns"
	pGcore "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/gcore"
	pGname "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/gname"
	pGoDaddy "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/godaddy"
	pHuaweiCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/huaweicloud"
	pNameDotCom "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namedotcom"
	pNameSilo "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/namesilo"
	pNS1 "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/ns1"
	pPowerDNS "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/powerdns"
	pRainYun "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/rainyun"
	pTencentCloud "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/tencentcloud"
	pVolcEngine "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/volcengine"
	pWestcn "github.com/usual2970/certimate/internal/pkg/core/applicant/acme-dns-01/lego-providers/westcn"
	"github.com/usual2970/certimate/internal/pkg/utils/maps"
)

func createApplicant(options *applicantOptions) (challenge.Provider, error) {
	/*
	  注意：如果追加新的常量值，请保持以 ASCII 排序。
	  NOTICE: If you add new constant, please keep ASCII order.
	*/
	switch options.Provider {
	case domain.ApplyDNSProviderTypeACMEHttpReq:
		{
			access := domain.AccessConfigForACMEHttpReq{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pACMEHttpReq.NewChallengeProvider(&pACMEHttpReq.ACMEHttpReqApplicantConfig{
				Endpoint:              access.Endpoint,
				Mode:                  access.Mode,
				Username:              access.Username,
				Password:              access.Password,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAliyun, domain.ApplyDNSProviderTypeAliyunDNS:
		{
			access := domain.AccessConfigForAliyun{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pAliyun.NewChallengeProvider(&pAliyun.AliyunApplicantConfig{
				AccessKeyId:           access.AccessKeyId,
				AccessKeySecret:       access.AccessKeySecret,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAWS, domain.ApplyDNSProviderTypeAWSRoute53:
		{
			access := domain.AccessConfigForAWS{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pAWSRoute53.NewChallengeProvider(&pAWSRoute53.AWSRoute53ApplicantConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				Region:                maps.GetValueAsString(options.ProviderApplyConfig, "region"),
				HostedZoneId:          maps.GetValueAsString(options.ProviderApplyConfig, "hostedZoneId"),
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeAzureDNS:
		{
			access := domain.AccessConfigForAzure{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pAzureDNS.NewChallengeProvider(&pAzureDNS.AzureDNSApplicantConfig{
				TenantId:              access.TenantId,
				ClientId:              access.ClientId,
				ClientSecret:          access.ClientSecret,
				CloudName:             access.CloudName,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeCloudflare:
		{
			access := domain.AccessConfigForCloudflare{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pCloudflare.NewChallengeProvider(&pCloudflare.CloudflareApplicantConfig{
				DnsApiToken:           access.DnsApiToken,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeClouDNS:
		{
			access := domain.AccessConfigForClouDNS{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pClouDNS.NewChallengeProvider(&pClouDNS.ClouDNSApplicantConfig{
				AuthId:                access.AuthId,
				AuthPassword:          access.AuthPassword,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeGcore:
		{
			access := domain.AccessConfigForGcore{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pGcore.NewChallengeProvider(&pGcore.GcoreApplicantConfig{
				ApiToken:              access.ApiToken,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeGname:
		{
			access := domain.AccessConfigForGname{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pGname.NewChallengeProvider(&pGname.GnameApplicantConfig{
				AppId:                 access.AppId,
				AppKey:                access.AppKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeGoDaddy:
		{
			access := domain.AccessConfigForGoDaddy{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pGoDaddy.NewChallengeProvider(&pGoDaddy.GoDaddyApplicantConfig{
				ApiKey:                access.ApiKey,
				ApiSecret:             access.ApiSecret,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeHuaweiCloud, domain.ApplyDNSProviderTypeHuaweiCloudDNS:
		{
			access := domain.AccessConfigForHuaweiCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pHuaweiCloud.NewChallengeProvider(&pHuaweiCloud.HuaweiCloudApplicantConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				Region:                maps.GetValueAsString(options.ProviderApplyConfig, "region"),
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNameDotCom:
		{
			access := domain.AccessConfigForNameDotCom{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pNameDotCom.NewChallengeProvider(&pNameDotCom.NameDotComApplicantConfig{
				Username:              access.Username,
				ApiToken:              access.ApiToken,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNameSilo:
		{
			access := domain.AccessConfigForNameSilo{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pNameSilo.NewChallengeProvider(&pNameSilo.NameSiloApplicantConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeNS1:
		{
			access := domain.AccessConfigForNS1{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pNS1.NewChallengeProvider(&pNS1.NS1ApplicantConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypePowerDNS:
		{
			access := domain.AccessConfigForPowerDNS{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pPowerDNS.NewChallengeProvider(&pPowerDNS.PowerDNSApplicantConfig{
				ApiUrl:                access.ApiUrl,
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeRainYun:
		{
			access := domain.AccessConfigForRainYun{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pRainYun.NewChallengeProvider(&pRainYun.RainYunApplicantConfig{
				ApiKey:                access.ApiKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeTencentCloud, domain.ApplyDNSProviderTypeTencentCloudDNS:
		{
			access := domain.AccessConfigForTencentCloud{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pTencentCloud.NewChallengeProvider(&pTencentCloud.TencentCloudApplicantConfig{
				SecretId:              access.SecretId,
				SecretKey:             access.SecretKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeVolcEngine, domain.ApplyDNSProviderTypeVolcEngineDNS:
		{
			access := domain.AccessConfigForVolcEngine{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pVolcEngine.NewChallengeProvider(&pVolcEngine.VolcEngineApplicantConfig{
				AccessKeyId:           access.AccessKeyId,
				SecretAccessKey:       access.SecretAccessKey,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}

	case domain.ApplyDNSProviderTypeWestcn:
		{
			access := domain.AccessConfigForWestcn{}
			if err := maps.Populate(options.ProviderAccessConfig, &access); err != nil {
				return nil, fmt.Errorf("failed to populate provider access config: %w", err)
			}

			applicant, err := pWestcn.NewChallengeProvider(&pWestcn.WestcnApplicantConfig{
				Username:              access.Username,
				ApiPassword:           access.ApiPassword,
				DnsPropagationTimeout: options.DnsPropagationTimeout,
				DnsTTL:                options.DnsTTL,
			})
			return applicant, err
		}
	}

	return nil, fmt.Errorf("unsupported applicant provider: %s", string(options.Provider))
}
