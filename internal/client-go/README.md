# Go API client for client

This is the API specification for Ory Identities with features such as registration, login, recovery, account verification, profile settings, password reset, identity management, session management, email and sms delivery, and more.


## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 
- Package version: 1.0.0
- Generator version: 7.12.0
- Build package: org.openapitools.codegen.languages.GoClientCodegen

## Installation

Install the following dependencies:

```sh
go get github.com/stretchr/testify/assert
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```go
import client "github.com/ory/client-go"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```go
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `client.ContextServerIndex` of type `int`.

```go
ctx := context.WithValue(context.Background(), client.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `client.ContextServerVariables` of type `map[string]string`.

```go
ctx := context.WithValue(context.Background(), client.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `client.ContextOperationServerIndices` and `client.ContextOperationServerVariables` context maps.

```go
ctx := context.WithValue(context.Background(), client.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), client.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*CourierAPI* | [**GetCourierMessage**](docs/CourierAPI.md#getcouriermessage) | **Get** /admin/courier/messages/{id} | Get a Message
*CourierAPI* | [**ListCourierMessages**](docs/CourierAPI.md#listcouriermessages) | **Get** /admin/courier/messages | List Messages
*FrontendAPI* | [**CreateBrowserLoginFlow**](docs/FrontendAPI.md#createbrowserloginflow) | **Get** /self-service/login/browser | Create Login Flow for Browsers
*FrontendAPI* | [**CreateBrowserLogoutFlow**](docs/FrontendAPI.md#createbrowserlogoutflow) | **Get** /self-service/logout/browser | Create a Logout URL for Browsers
*FrontendAPI* | [**CreateBrowserRecoveryFlow**](docs/FrontendAPI.md#createbrowserrecoveryflow) | **Get** /self-service/recovery/browser | Create Recovery Flow for Browsers
*FrontendAPI* | [**CreateBrowserRegistrationFlow**](docs/FrontendAPI.md#createbrowserregistrationflow) | **Get** /self-service/registration/browser | Create Registration Flow for Browsers
*FrontendAPI* | [**CreateBrowserSettingsFlow**](docs/FrontendAPI.md#createbrowsersettingsflow) | **Get** /self-service/settings/browser | Create Settings Flow for Browsers
*FrontendAPI* | [**CreateBrowserVerificationFlow**](docs/FrontendAPI.md#createbrowserverificationflow) | **Get** /self-service/verification/browser | Create Verification Flow for Browser Clients
*FrontendAPI* | [**CreateFedcmFlow**](docs/FrontendAPI.md#createfedcmflow) | **Get** /self-service/fed-cm/parameters | Get FedCM Parameters
*FrontendAPI* | [**CreateNativeLoginFlow**](docs/FrontendAPI.md#createnativeloginflow) | **Get** /self-service/login/api | Create Login Flow for Native Apps
*FrontendAPI* | [**CreateNativeRecoveryFlow**](docs/FrontendAPI.md#createnativerecoveryflow) | **Get** /self-service/recovery/api | Create Recovery Flow for Native Apps
*FrontendAPI* | [**CreateNativeRegistrationFlow**](docs/FrontendAPI.md#createnativeregistrationflow) | **Get** /self-service/registration/api | Create Registration Flow for Native Apps
*FrontendAPI* | [**CreateNativeSettingsFlow**](docs/FrontendAPI.md#createnativesettingsflow) | **Get** /self-service/settings/api | Create Settings Flow for Native Apps
*FrontendAPI* | [**CreateNativeVerificationFlow**](docs/FrontendAPI.md#createnativeverificationflow) | **Get** /self-service/verification/api | Create Verification Flow for Native Apps
*FrontendAPI* | [**DisableMyOtherSessions**](docs/FrontendAPI.md#disablemyothersessions) | **Delete** /sessions | Disable my other sessions
*FrontendAPI* | [**DisableMySession**](docs/FrontendAPI.md#disablemysession) | **Delete** /sessions/{id} | Disable one of my sessions
*FrontendAPI* | [**ExchangeSessionToken**](docs/FrontendAPI.md#exchangesessiontoken) | **Get** /sessions/token-exchange | Exchange Session Token
*FrontendAPI* | [**GetFlowError**](docs/FrontendAPI.md#getflowerror) | **Get** /self-service/errors | Get User-Flow Errors
*FrontendAPI* | [**GetLoginFlow**](docs/FrontendAPI.md#getloginflow) | **Get** /self-service/login/flows | Get Login Flow
*FrontendAPI* | [**GetRecoveryFlow**](docs/FrontendAPI.md#getrecoveryflow) | **Get** /self-service/recovery/flows | Get Recovery Flow
*FrontendAPI* | [**GetRegistrationFlow**](docs/FrontendAPI.md#getregistrationflow) | **Get** /self-service/registration/flows | Get Registration Flow
*FrontendAPI* | [**GetSettingsFlow**](docs/FrontendAPI.md#getsettingsflow) | **Get** /self-service/settings/flows | Get Settings Flow
*FrontendAPI* | [**GetVerificationFlow**](docs/FrontendAPI.md#getverificationflow) | **Get** /self-service/verification/flows | Get Verification Flow
*FrontendAPI* | [**GetWebAuthnJavaScript**](docs/FrontendAPI.md#getwebauthnjavascript) | **Get** /.well-known/ory/webauthn.js | Get WebAuthn JavaScript
*FrontendAPI* | [**ListMySessions**](docs/FrontendAPI.md#listmysessions) | **Get** /sessions | Get My Active Sessions
*FrontendAPI* | [**PerformNativeLogout**](docs/FrontendAPI.md#performnativelogout) | **Delete** /self-service/logout/api | Perform Logout for Native Apps
*FrontendAPI* | [**ToSession**](docs/FrontendAPI.md#tosession) | **Get** /sessions/whoami | Check Who the Current HTTP Session Belongs To
*FrontendAPI* | [**UpdateFedcmFlow**](docs/FrontendAPI.md#updatefedcmflow) | **Post** /self-service/fed-cm/token | Submit a FedCM token
*FrontendAPI* | [**UpdateLoginFlow**](docs/FrontendAPI.md#updateloginflow) | **Post** /self-service/login | Submit a Login Flow
*FrontendAPI* | [**UpdateLogoutFlow**](docs/FrontendAPI.md#updatelogoutflow) | **Get** /self-service/logout | Update Logout Flow
*FrontendAPI* | [**UpdateRecoveryFlow**](docs/FrontendAPI.md#updaterecoveryflow) | **Post** /self-service/recovery | Update Recovery Flow
*FrontendAPI* | [**UpdateRegistrationFlow**](docs/FrontendAPI.md#updateregistrationflow) | **Post** /self-service/registration | Update Registration Flow
*FrontendAPI* | [**UpdateSettingsFlow**](docs/FrontendAPI.md#updatesettingsflow) | **Post** /self-service/settings | Complete Settings Flow
*FrontendAPI* | [**UpdateVerificationFlow**](docs/FrontendAPI.md#updateverificationflow) | **Post** /self-service/verification | Complete Verification Flow
*IdentityAPI* | [**BatchPatchIdentities**](docs/IdentityAPI.md#batchpatchidentities) | **Patch** /admin/identities | Create multiple identities
*IdentityAPI* | [**CreateIdentity**](docs/IdentityAPI.md#createidentity) | **Post** /admin/identities | Create an Identity
*IdentityAPI* | [**CreateRecoveryCodeForIdentity**](docs/IdentityAPI.md#createrecoverycodeforidentity) | **Post** /admin/recovery/code | Create a Recovery Code
*IdentityAPI* | [**CreateRecoveryLinkForIdentity**](docs/IdentityAPI.md#createrecoverylinkforidentity) | **Post** /admin/recovery/link | Create a Recovery Link
*IdentityAPI* | [**DeleteIdentity**](docs/IdentityAPI.md#deleteidentity) | **Delete** /admin/identities/{id} | Delete an Identity
*IdentityAPI* | [**DeleteIdentityCredentials**](docs/IdentityAPI.md#deleteidentitycredentials) | **Delete** /admin/identities/{id}/credentials/{type} | Delete a credential for a specific identity
*IdentityAPI* | [**DeleteIdentitySessions**](docs/IdentityAPI.md#deleteidentitysessions) | **Delete** /admin/identities/{id}/sessions | Delete &amp; Invalidate an Identity&#39;s Sessions
*IdentityAPI* | [**DisableSession**](docs/IdentityAPI.md#disablesession) | **Delete** /admin/sessions/{id} | Deactivate a Session
*IdentityAPI* | [**ExtendSession**](docs/IdentityAPI.md#extendsession) | **Patch** /admin/sessions/{id}/extend | Extend a Session
*IdentityAPI* | [**GetIdentity**](docs/IdentityAPI.md#getidentity) | **Get** /admin/identities/{id} | Get an Identity
*IdentityAPI* | [**GetIdentityByExternalID**](docs/IdentityAPI.md#getidentitybyexternalid) | **Get** /admin/identities/by/external/{externalID} | Get an Identity by its External ID
*IdentityAPI* | [**GetIdentitySchema**](docs/IdentityAPI.md#getidentityschema) | **Get** /schemas/{id} | Get Identity JSON Schema
*IdentityAPI* | [**GetSession**](docs/IdentityAPI.md#getsession) | **Get** /admin/sessions/{id} | Get Session
*IdentityAPI* | [**ListIdentities**](docs/IdentityAPI.md#listidentities) | **Get** /admin/identities | List Identities
*IdentityAPI* | [**ListIdentitySchemas**](docs/IdentityAPI.md#listidentityschemas) | **Get** /schemas | Get all Identity Schemas
*IdentityAPI* | [**ListIdentitySessions**](docs/IdentityAPI.md#listidentitysessions) | **Get** /admin/identities/{id}/sessions | List an Identity&#39;s Sessions
*IdentityAPI* | [**ListSessions**](docs/IdentityAPI.md#listsessions) | **Get** /admin/sessions | List All Sessions
*IdentityAPI* | [**PatchIdentity**](docs/IdentityAPI.md#patchidentity) | **Patch** /admin/identities/{id} | Patch an Identity
*IdentityAPI* | [**UpdateIdentity**](docs/IdentityAPI.md#updateidentity) | **Put** /admin/identities/{id} | Update an Identity
*MetadataAPI* | [**GetVersion**](docs/MetadataAPI.md#getversion) | **Get** /version | Return Running Software Version.
*MetadataAPI* | [**IsAlive**](docs/MetadataAPI.md#isalive) | **Get** /health/alive | Check HTTP Server Status
*MetadataAPI* | [**IsReady**](docs/MetadataAPI.md#isready) | **Get** /health/ready | Check HTTP Server and Database Status


## Documentation For Models

 - [AuthenticatorAssuranceLevel](docs/AuthenticatorAssuranceLevel.md)
 - [BatchPatchIdentitiesResponse](docs/BatchPatchIdentitiesResponse.md)
 - [ConsistencyRequestParameters](docs/ConsistencyRequestParameters.md)
 - [ContinueWith](docs/ContinueWith.md)
 - [ContinueWithRecoveryUi](docs/ContinueWithRecoveryUi.md)
 - [ContinueWithRecoveryUiFlow](docs/ContinueWithRecoveryUiFlow.md)
 - [ContinueWithRedirectBrowserTo](docs/ContinueWithRedirectBrowserTo.md)
 - [ContinueWithSetOrySessionToken](docs/ContinueWithSetOrySessionToken.md)
 - [ContinueWithSettingsUi](docs/ContinueWithSettingsUi.md)
 - [ContinueWithSettingsUiFlow](docs/ContinueWithSettingsUiFlow.md)
 - [ContinueWithVerificationUi](docs/ContinueWithVerificationUi.md)
 - [ContinueWithVerificationUiFlow](docs/ContinueWithVerificationUiFlow.md)
 - [CourierMessageStatus](docs/CourierMessageStatus.md)
 - [CourierMessageType](docs/CourierMessageType.md)
 - [CreateFedcmFlowResponse](docs/CreateFedcmFlowResponse.md)
 - [CreateIdentityBody](docs/CreateIdentityBody.md)
 - [CreateRecoveryCodeForIdentityBody](docs/CreateRecoveryCodeForIdentityBody.md)
 - [CreateRecoveryLinkForIdentityBody](docs/CreateRecoveryLinkForIdentityBody.md)
 - [DeleteMySessionsCount](docs/DeleteMySessionsCount.md)
 - [ErrorAuthenticatorAssuranceLevelNotSatisfied](docs/ErrorAuthenticatorAssuranceLevelNotSatisfied.md)
 - [ErrorBrowserLocationChangeRequired](docs/ErrorBrowserLocationChangeRequired.md)
 - [ErrorFlowReplaced](docs/ErrorFlowReplaced.md)
 - [ErrorGeneric](docs/ErrorGeneric.md)
 - [FlowError](docs/FlowError.md)
 - [GenericError](docs/GenericError.md)
 - [GetVersion200Response](docs/GetVersion200Response.md)
 - [HealthNotReadyStatus](docs/HealthNotReadyStatus.md)
 - [HealthStatus](docs/HealthStatus.md)
 - [Identity](docs/Identity.md)
 - [IdentityCredentials](docs/IdentityCredentials.md)
 - [IdentityCredentialsCode](docs/IdentityCredentialsCode.md)
 - [IdentityCredentialsCodeAddress](docs/IdentityCredentialsCodeAddress.md)
 - [IdentityCredentialsOidc](docs/IdentityCredentialsOidc.md)
 - [IdentityCredentialsOidcProvider](docs/IdentityCredentialsOidcProvider.md)
 - [IdentityCredentialsPassword](docs/IdentityCredentialsPassword.md)
 - [IdentityPatch](docs/IdentityPatch.md)
 - [IdentityPatchResponse](docs/IdentityPatchResponse.md)
 - [IdentitySchemaContainer](docs/IdentitySchemaContainer.md)
 - [IdentityWithCredentials](docs/IdentityWithCredentials.md)
 - [IdentityWithCredentialsOidc](docs/IdentityWithCredentialsOidc.md)
 - [IdentityWithCredentialsOidcConfig](docs/IdentityWithCredentialsOidcConfig.md)
 - [IdentityWithCredentialsOidcConfigProvider](docs/IdentityWithCredentialsOidcConfigProvider.md)
 - [IdentityWithCredentialsPassword](docs/IdentityWithCredentialsPassword.md)
 - [IdentityWithCredentialsPasswordConfig](docs/IdentityWithCredentialsPasswordConfig.md)
 - [IdentityWithCredentialsSaml](docs/IdentityWithCredentialsSaml.md)
 - [IdentityWithCredentialsSamlConfig](docs/IdentityWithCredentialsSamlConfig.md)
 - [IdentityWithCredentialsSamlConfigProvider](docs/IdentityWithCredentialsSamlConfigProvider.md)
 - [IsAlive200Response](docs/IsAlive200Response.md)
 - [IsReady503Response](docs/IsReady503Response.md)
 - [JsonPatch](docs/JsonPatch.md)
 - [LoginFlow](docs/LoginFlow.md)
 - [LoginFlowState](docs/LoginFlowState.md)
 - [LogoutFlow](docs/LogoutFlow.md)
 - [Message](docs/Message.md)
 - [MessageDispatch](docs/MessageDispatch.md)
 - [NeedsPrivilegedSessionError](docs/NeedsPrivilegedSessionError.md)
 - [OAuth2Client](docs/OAuth2Client.md)
 - [OAuth2ConsentRequestOpenIDConnectContext](docs/OAuth2ConsentRequestOpenIDConnectContext.md)
 - [OAuth2LoginRequest](docs/OAuth2LoginRequest.md)
 - [PatchIdentitiesBody](docs/PatchIdentitiesBody.md)
 - [PerformNativeLogoutBody](docs/PerformNativeLogoutBody.md)
 - [Provider](docs/Provider.md)
 - [RecoveryCodeForIdentity](docs/RecoveryCodeForIdentity.md)
 - [RecoveryFlow](docs/RecoveryFlow.md)
 - [RecoveryFlowState](docs/RecoveryFlowState.md)
 - [RecoveryIdentityAddress](docs/RecoveryIdentityAddress.md)
 - [RecoveryLinkForIdentity](docs/RecoveryLinkForIdentity.md)
 - [RegistrationFlow](docs/RegistrationFlow.md)
 - [RegistrationFlowState](docs/RegistrationFlowState.md)
 - [SelfServiceFlowExpiredError](docs/SelfServiceFlowExpiredError.md)
 - [Session](docs/Session.md)
 - [SessionAuthenticationMethod](docs/SessionAuthenticationMethod.md)
 - [SessionDevice](docs/SessionDevice.md)
 - [SettingsFlow](docs/SettingsFlow.md)
 - [SettingsFlowState](docs/SettingsFlowState.md)
 - [SuccessfulCodeExchangeResponse](docs/SuccessfulCodeExchangeResponse.md)
 - [SuccessfulNativeLogin](docs/SuccessfulNativeLogin.md)
 - [SuccessfulNativeRegistration](docs/SuccessfulNativeRegistration.md)
 - [TokenPagination](docs/TokenPagination.md)
 - [TokenPaginationHeaders](docs/TokenPaginationHeaders.md)
 - [UiContainer](docs/UiContainer.md)
 - [UiNode](docs/UiNode.md)
 - [UiNodeAnchorAttributes](docs/UiNodeAnchorAttributes.md)
 - [UiNodeAttributes](docs/UiNodeAttributes.md)
 - [UiNodeDivisionAttributes](docs/UiNodeDivisionAttributes.md)
 - [UiNodeImageAttributes](docs/UiNodeImageAttributes.md)
 - [UiNodeInputAttributes](docs/UiNodeInputAttributes.md)
 - [UiNodeMeta](docs/UiNodeMeta.md)
 - [UiNodeScriptAttributes](docs/UiNodeScriptAttributes.md)
 - [UiNodeTextAttributes](docs/UiNodeTextAttributes.md)
 - [UiText](docs/UiText.md)
 - [UpdateFedcmFlowBody](docs/UpdateFedcmFlowBody.md)
 - [UpdateIdentityBody](docs/UpdateIdentityBody.md)
 - [UpdateLoginFlowBody](docs/UpdateLoginFlowBody.md)
 - [UpdateLoginFlowWithCodeMethod](docs/UpdateLoginFlowWithCodeMethod.md)
 - [UpdateLoginFlowWithIdentifierFirstMethod](docs/UpdateLoginFlowWithIdentifierFirstMethod.md)
 - [UpdateLoginFlowWithLookupSecretMethod](docs/UpdateLoginFlowWithLookupSecretMethod.md)
 - [UpdateLoginFlowWithOidcMethod](docs/UpdateLoginFlowWithOidcMethod.md)
 - [UpdateLoginFlowWithPasskeyMethod](docs/UpdateLoginFlowWithPasskeyMethod.md)
 - [UpdateLoginFlowWithPasswordMethod](docs/UpdateLoginFlowWithPasswordMethod.md)
 - [UpdateLoginFlowWithSamlMethod](docs/UpdateLoginFlowWithSamlMethod.md)
 - [UpdateLoginFlowWithTotpMethod](docs/UpdateLoginFlowWithTotpMethod.md)
 - [UpdateLoginFlowWithWebAuthnMethod](docs/UpdateLoginFlowWithWebAuthnMethod.md)
 - [UpdateRecoveryFlowBody](docs/UpdateRecoveryFlowBody.md)
 - [UpdateRecoveryFlowWithCodeMethod](docs/UpdateRecoveryFlowWithCodeMethod.md)
 - [UpdateRecoveryFlowWithLinkMethod](docs/UpdateRecoveryFlowWithLinkMethod.md)
 - [UpdateRegistrationFlowBody](docs/UpdateRegistrationFlowBody.md)
 - [UpdateRegistrationFlowWithCodeMethod](docs/UpdateRegistrationFlowWithCodeMethod.md)
 - [UpdateRegistrationFlowWithOidcMethod](docs/UpdateRegistrationFlowWithOidcMethod.md)
 - [UpdateRegistrationFlowWithPasskeyMethod](docs/UpdateRegistrationFlowWithPasskeyMethod.md)
 - [UpdateRegistrationFlowWithPasswordMethod](docs/UpdateRegistrationFlowWithPasswordMethod.md)
 - [UpdateRegistrationFlowWithProfileMethod](docs/UpdateRegistrationFlowWithProfileMethod.md)
 - [UpdateRegistrationFlowWithSamlMethod](docs/UpdateRegistrationFlowWithSamlMethod.md)
 - [UpdateRegistrationFlowWithWebAuthnMethod](docs/UpdateRegistrationFlowWithWebAuthnMethod.md)
 - [UpdateSettingsFlowBody](docs/UpdateSettingsFlowBody.md)
 - [UpdateSettingsFlowWithLookupMethod](docs/UpdateSettingsFlowWithLookupMethod.md)
 - [UpdateSettingsFlowWithOidcMethod](docs/UpdateSettingsFlowWithOidcMethod.md)
 - [UpdateSettingsFlowWithPasskeyMethod](docs/UpdateSettingsFlowWithPasskeyMethod.md)
 - [UpdateSettingsFlowWithPasswordMethod](docs/UpdateSettingsFlowWithPasswordMethod.md)
 - [UpdateSettingsFlowWithProfileMethod](docs/UpdateSettingsFlowWithProfileMethod.md)
 - [UpdateSettingsFlowWithSamlMethod](docs/UpdateSettingsFlowWithSamlMethod.md)
 - [UpdateSettingsFlowWithTotpMethod](docs/UpdateSettingsFlowWithTotpMethod.md)
 - [UpdateSettingsFlowWithWebAuthnMethod](docs/UpdateSettingsFlowWithWebAuthnMethod.md)
 - [UpdateVerificationFlowBody](docs/UpdateVerificationFlowBody.md)
 - [UpdateVerificationFlowWithCodeMethod](docs/UpdateVerificationFlowWithCodeMethod.md)
 - [UpdateVerificationFlowWithLinkMethod](docs/UpdateVerificationFlowWithLinkMethod.md)
 - [VerifiableIdentityAddress](docs/VerifiableIdentityAddress.md)
 - [VerificationFlow](docs/VerificationFlow.md)
 - [VerificationFlowState](docs/VerificationFlowState.md)
 - [Version](docs/Version.md)


## Documentation For Authorization


Authentication schemes defined for the API:
### oryAccessToken

- **Type**: API key
- **API key parameter name**: Authorization
- **Location**: HTTP header

Note, each API key must be added to a map of `map[string]APIKey` where the key is: oryAccessToken and passed in as the auth context for each request.

Example

```go
auth := context.WithValue(
		context.Background(),
		client.ContextAPIKeys,
		map[string]client.APIKey{
			"oryAccessToken": {Key: "API_KEY_STRING"},
		},
	)
r, err := client.Service.Operation(auth, args)
```


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author

office@ory.sh

