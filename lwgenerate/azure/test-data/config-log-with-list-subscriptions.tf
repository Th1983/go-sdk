terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.16"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 2.91.0"
    }
  }
}

provider "azuread" {
}

provider "azurerm" {
  features {
  }
}

module "az_ad_application" {
  source  = "lacework/ad-application/azure"
  version = "~> 1.0"
}

module "az_config" {
  source                      = "lacework/config/azure"
  version                     = "~> 1.0"
  application_id              = "module.az_ad_application.application_id"
  application_password        = "module.az_ad_application.application_password"
  service_principal_id        = "module.az_ad_application.service_principal_id"
  subscription_ids            = ["test-id-1", "test-id-2", "test-id-3"]
  use_existing_ad_application = true
}