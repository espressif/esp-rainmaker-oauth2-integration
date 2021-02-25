## ESP Rainmaker Oauth2 Integration

Following are the steps to deploy ESP Rainmaker Ouath2 Integration respository in your account:

1. In the Application settings you have to enter the following URLs
      - RainmakerOauth2AuthorizeUrl : The authorize url for your IDP to authorize user
      - RainmakerOauth2EmailUrl : The email url for your IDP to fetch user email details (optional)
      - RainmakerOauth2TokenUrl : The token url for your IDP to fetch the user authentication tokens
      - RainmakerOauth2UserInfoUrl : The userinfo url for your IDP to fetch the user details

**Note: If at the time of deploying this OAuth2 integration repository, if you do not have above urls, we can still proceed with the deployment. These URLs can be configured later on, using the configuration APIs provided with this repository.**
   
2. Click on the checkbox - “I acknowledge that this app creates custom IAM roles”.
3. Click on Deploy button

This will trigger the deployment of ESP Rainmaker Oauth2 Integration in your account.
