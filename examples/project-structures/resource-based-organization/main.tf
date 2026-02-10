

module "users" {
  source = "./modules/01-users"

  # input variables
  user1 = "bob@company.io"
  user2 = "tim@company.io"

  # provider configuration
  providers = {
    conduktor = conduktor.console
  }
}

module "groups" {
  source = "./modules/02-groups"

  # input variables
  group_name = "website-analytics-team"
  users      = module.users.users_list

  # provider configuration
  providers = {
    conduktor = conduktor.console
  }
}

module "clusters" {
  source = "./modules/03-clusters"

  # input variables
  cluster_name = "my-cluster"

  providers = {
    conduktor = conduktor.console
  }
}

module "interceptors" {
  source = "./modules/04-interceptors"

  # provider configuration
  providers = {
    conduktor = conduktor.gateway
  }
}

module "topic-policies" {
  source = "./modules/05-topic-policies"

  # provider configuration
  providers = {
    conduktor = conduktor.console
  }
}

module "applications" {
  source = "./modules/06-applications"

  # input variables
  application_name = "website-analytics"
  title            = "Website Analytics"
  description      = "Application for streaming web analytics"
  owner            = "website-analytics-team"

  # provider configuration
  providers = {
    conduktor = conduktor.console
  }
}

module "application-instances" {
  source = "./modules/07-application-instances"

  # provider configuration
  providers = {
    conduktor = conduktor.console
  }
}

# In this example topics are being created from this root, centrally.
# Topics meet topic policy criteria
module "topics" {
  source = "./modules/08-topics-a"

  # provider configuration
  providers = {
    conduktor = conduktor.console
  }
}

# # In this example topic creation is delegated, an appInstance token is used by the app team 
# # Apply the topic resource from the module directly
# # Topics deliberately fail to meet topic policy criteria and will need adjusting based on error feedback
# module "topics" {
#   source = "./modules/08-topics-b"

#   # provider configuration
#   providers = {
#     conduktor = conduktor.console
#   }
# }