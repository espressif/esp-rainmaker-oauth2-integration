## Espressif RainMaker OAuth2 Integration

This project contains the required modules for integration with third-party OAuth2 providers.

### What is OAuth2 Integration?

OAuth2 is a protocol used for Authorization. OAuth2 app can be integrated with OpenId Connect, allowing it to be used for user authentication. 

You can enable your users who already have accounts with OpenID Connect (OIDC) identity providers (IdPs) (like Github, Salesforce, or Ping Identity) to skip the sign-up step and sign in to your application using an existing account.

To facilitate Sign-in with third-party IDPs, RainMaker provides support for OAuth2 integration with IDPs.

You can configure the required Identity Provider with help of the example given in this repo.

### How does OAuth2 Integration work?
The Oauth2 integration includes API endpoints, which need to be configured in the AWS account, where RainMaker is deployed.
The API endpoints are called by AWS Cognito, to sign-in users.

When a user clicks on the sign-in button on RainMaker dashboard, AWS Cognito makes a call to `esp-Oauth2Authorize` lambda. This lambda in turn makes a GET API call to your IDP's authorize API endpoint(specified as env variable Rainmaker_Oauth2_Authorize_URL, espoauth2integration.yml file)
The query parameters that are required to make Authorization successful, are added to the URL. In the response of this URL, we get a code and state. 

The application requests authorization to access service resources from the user. If the user authorizes the request, the application receives an authorization grant, which is `code` for grant_type=authorization_code. 

AWS Cognito then makes a call to `esp-Oauth2Token` lambda, which in turn requests an access token from IDP's token endpoint(specified as env variable Rainmaker_Oauth2_Token_URL), by passing the authorization code along with authentication details, including the client secret, to the API token endpoint. 
In the response, Cognito gets the access_token, which has limited access to user's resources in IDP. Here we have given access to read the user's email and email_verified attributes.

Using the access token received with the `esp-Oauth2Token` endpoint, Cognito fetches user info and passes this info to RainMaker.

### Deploying and Configuring the OAuth2 Integration consist of the following high-level steps:
  - [Deploying the SAR repository from the AWS Console](#deploying-rainmaker-oauth2-integration-sar-repository-from-the-aws-console)
  - [Registering an OAuth2 App with the Identity Provider](#registering-an-oauth2-app-with-the-identity-provider)
  - [Configuring the Domain Name in AWS Cognito](#configuring-the-domain-name-in-aws-cognito)
  - [Configuring the IDP in AWS Cognito](#configuring-the-idp-in-aws-cognito)
  - [Configuring the Cognito App Client Settings](#configuring-the-cognito-app-client-settings)

### Deploying RainMaker OAuth2 Integration SAR repository from the AWS Console
#### Pre-requisites for starting RainMaker OAuth2 Integration
 - ESP RainMaker should be deployed into your AWS account.
 - Authentication URLs for your IDP like authorize url, token url, userinfo url etc. (optional)
   
#### Deployment Steps
1. Log in to AWS Management Console.
2. In the Services, search for SAR (Serverless Application Repository) and click on SAR to navigate to SAR service, as shown below.
![](./images/search_SAR.png)

3. Click on `Available applications` and select the `Private applications` tab. Click the checkbox - Show apps that create custom IAM roles or resource policies and search for `espoauth2integration`  as shown below:  

If you don’t see the espoauth2integration repository, please make sure you are in the correct region, else contact RainMaker Support Team. (esp-rainmaker-support@espressif.com)
![s](./images/search_oauth2_repo_sar.png)  

4. Click on the `espoauth2integration` repository.  
   
5. After clicking on `espoauth2integration` SAR repository, the details for the repository will be displayed as below.

![](./images/oauth2_sar_repo.png)
   
6. Scroll down and you will see the `Application settings` as shown below.

7. In the Application settings you have to enter the following URLs
- RainmakerOauth2AuthorizeUrl : The authorize url for your IDP to authorize the user
- RainmakerOauth2EmailUrl : The email url for your IDP to fetch user email details (optional)
- RainmakerOauth2TokenUrl : The token url for your IDP to fetch the user authentication tokens
- RainmakerOauth2UserInfoUrl : The userinfo url for your IDP to fetch the user details

<div align="center"><img src = "./images/edit_urls_sar.png" width = 400 height = 500 ></div>

   **Note: If at the time of deploying this OAuth2 integration repository, you do not have the above URLs, we can still proceed with the deployment. These URLs can be configured later on, using the configuration APIs provided with this repository.**

8. Click on the checkbox - “I acknowledge that this app creates custom IAM roles”.
9. Click on the Deploy button, the deployment will start and the below screen will be displayed. The deployment process will be completed in a few minutes.
![](./images/deployment_screen.png)

10. To check the deployment status, click on the `Deployments` tab. After successful deployment, the status should look like the following.
![](./images/deployment_status.png)


### Registering an OAuth2 App with the Identity Provider

After deployment is successful, the first step is to register an OAuth2 App with the Identity provider, you are going to integrate with. This is generally done on the IDP's developer portal.

While registering the app, you will require the Callback URL during the configuration steps. To get this Callback URL, follow the below steps:
1. Log in to AWS Management Console.
2. In the Services, search for Cloudformation and click on it to navigate to the Cloudformation console, as shown below.
![](./images/search_cfd.png)
3. Toggle the `view nested` button to disable it and search for `OAuth2`, as shown below.
![](./images/search_oauth2.png)
4. Click on `serverlessrepo-ESP-Rainmaker-OAuth2-Integration`.
5. Click on the `Outputs` tab and search for `Oauth2CognitoCallbackURL` as shown below.
![](./images/search_callback_url.png)
6. In the Outputs section, note down the `Oauth2CognitoCallbackURL` value. This will be used as Callback URL.

Use the above callback URL and register your app. After successful registration, you will get ClientId and Client Secret Key. Every IDP provides authentication urls like authorize url, token url, userinfo url etc. You need to note down these URLs which are used at the time of deployment. The URLs can be found on the docs provided by the identity providers.

Following is an example link to get the Authentication URL, when Salesforce is used as Identity Provider:
- Salesforce - https://help.salesforce.com/articleView?id=sf.remoteaccess_oauth_endpoints.htm&type=5

**Note :** Some IDPs call ClientId as AppId and Client Secret Key as App Secret Key. In this documentation, we are using the terms ClientId and Client Secret Key.

### Configuring the Domain Name in AWS Cognito

Steps to configure the domain name in ESP RainMaker Cognito userpool are shown below.
1. Log in to AWS Management Console.
2. In the Services, search for Cognito and click on Cognito to navigate to Cognito service, as shown below.
![](./images/search_cognito.png)
3. On the Cognito console, click on `Manage User Pools` and select `rainmaker-user-pool` from the available user pools.
4. From the left pane, select `App integration->Domain name` as shown below.

<div align="center"><img src = "./images/select_domain_name.png" width = 400 height = 200 ></div>

5. Type the domain prefix you want to use in the Domain prefix.
![](./images/enter_domain_name.png)

6. Choose Check availability to confirm the domain prefix is available. If the domain name is available it will show you the below success message.
![](./images/available_domain_name.png)

7. Choose to save changes.

### Configuring the IDP in AWS Cognito

Steps to configure identity provider details in ESP RainMaker Cognito userpool are as shown below.
1. Log in to AWS Management Console.
2. In the Services, search for Cognito and click on Cognito to navigate to Cognito service, as shown below.
![](./images/search_cognito.png)
3. On the Cognito console, click on `Manage User Pools` and select `rainmaker-user-pool` from the available user pools.
4. From the left pane, select `Federation->Identity providers` as shown below.

<div align="center"><img src = "./images/click_identity_provider.png" width = 300 height = 500 ></div>

5. On the Identity providers page, select `OpenID Connect`, as shown below
![](./images/oidc.png)
 * 5.1. For the Provider name field, add the name of your IDP [eg. GitHub, Salesforce, etc]
 * 5.2. For the Client ID field, add the Client ID that you have got after registration of your app with IDP.
 * 5.3. For the Client secret field, add the Client Secret Key that you have got after registration of your app with IDP.
 * 5.4. For the `authorize_scope' field, specify the details that you want to fetch for users at the time of sign-in. You can check IDP's development docs to get available options for scope.
 * 5.5. For the 'issuer' field, enter the endpoint, which you get in the output section while deploying the espoauth2integration template. It is the base path of your RainMaker APIs.
 * 5.6. Click on the `Run Discovery` button. If Discovery returned no results, manually add the fields for the endpoints. To get the values for this endpoint go to the output section of espoauth2integration in the Cloudformation console.
 The values for endpoints are as follow:
   * Authorization endpoint - {{Oauth2IssuerApiResource}}/oauth2_authorize
   * Token endpoint - {{Oauth2IssuerApiResource}}/oauth2_token
   * UserInfo endpoint - {{Oauth2IssuerApiResource}}/oauth2_userinfo
   * JWKS endpoint - {{Oauth2IssuerApiResource}}/oauth2_jwks
6. Click on create provider.

 #### Configuring the Cognito App Client Settings
 
1. Log in to AWS Management Console.
2. In the Services, search for Cognito and click on Cognito to navigate to Cognito service, as shown below.
![](./images/search_cognito.png)
3. On the Cognito console, click on `Manage User Pools` and select `rainmaker-user-pool` from the available user pools. 
4. On the Cognito console, from the left pane click on `App client Settings`.

<div align="center"><img src = "./images/select_app_client.png" width = 300 height = 500 ></div>

5. Go to App client `rainmaker-user-pool-client` as shown below
![](./images/rmaker_app_client.png)
6. In the `Enabled Identity Providers` section, select all required IDPs. The Identity Provider which is configured should be selected.
7. In the `Callback URL(s)` section enter the callback URL where you want to redirect the user after successful login.
8. Click on the 'save changes' button.

### Steps to verify IDP integration with RainMaker:

1. Log in to AWS Management Console.
2. In the Services, search for Cognito and click on Cognito to navigate to Cognito service, as shown below.
![](./images/search_cognito.png)
3. On the Cognito console, click on `Manage User Pools` and select `rainmaker-user-pool` from the available user pools. 
4. On the Cognito console, from the left pane click on `App client Settings`.

<div align="center"><img src = "./images/select_app_client.png" width = 300 height = 500 ></div>

5. Go to App client `rainmaker-user-pool-client` and click on `Launch Hosted UI` as shown below.

<div align="center"><img src = "./images/hosted_ui.png" width = 800 height = 200></div>

6. You will be redirected to a login window. This window will show you options to log in with the IDP that you have configured. For example, if you are using GitHub as the Identity Provider, below login screen will be displayed:

<div align="center"><img src = "./images/GitHub_login.png" width = 400 height = 250 ></div>

7. Click on GitHub, and you will be redirected to the GitHub login console. Enter the username and password to log in.
8. If all the above integrations are correct, login should be successful and you will be redirected to the callback URL that you have provided in the last section.