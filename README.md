# Identity Provider based on JWT's 🔐

I needed a small authentication service, so I created this service that handles user authentication using JSON Web Tokens (JWTs). Let's dive into the details! 🚀

## Certificate Storage 📜

To ensure a secure and efficient way to store certificates on non-persistent systems, this service reads the certificates from MongoDB.

Before getting started, please refer to the `.env.example` file to configure your service according to your specific requirements.

## Build and Run 🏃‍♀️

You can build the service by executing the following command:
```
go build -o ./build ./src/main.go
```

To run the service, use the following command:
```
./build
```

## Routes 🛣️

### Login 🔑

To initiate the login process, navigate to `/login?provider=<provider>`. Replace `<provider>` with the desired authentication provider such as GitHub, Google, or Discord. This will redirect you to the OAuth screen of the selected provider, where you can authenticate yourself securely.

### Callbacks 🔄

For each authentication provider, you need to add a callback URL. The callback URL should follow this format: `/callback/provider`, where `provider` corresponds to the authentication provider you are integrating (e.g., `/callback/google` for Google authentication). After successful authentication, the provider will redirect the user back to the specified callback URL.

### User Info 👤

To retrieve user information for a specific token, you must store the token as a cookie named `access_token`. Once the token cookie is set, you can navigate to `/api/user` to fetch the user information associated with the token. This endpoint provides a convenient way to access user details after successful authentication.

### User Sessions 📆

This service allows you to manage user sessions effectively. You can view all the tokens (sessions) that are currently valid for a particular user. Additionally, you have the option to blacklist specific tokens by blocking their corresponding token ID.

To see all the tokens/sessions for a user, simply navigate to `/api/user/sessions/`. This endpoint provides an overview of all the active sessions associated with the user.

If you need to validate whether a token has been revoked, visit `/api/user/sessions/<tokenid>`. This endpoint will respond with information about whether the token has been revoked, allowing you to maintain control over the validity of user sessions.

To revoke/invalidate a token/session, send a `DELETE` method to `/api/user/sessions/<tokenid>`. This endpoint will invalidate
the token with the assoiciated `tokenid` and will block future requests with that token. When invalidating your app can also
store the `tokenid`, so it cuts the round trip for checking if the token with `tokenid` is valid.

To revoke/invalidate all tokens/sessions, send a `DELETE` method to `/api/users/sessions/invalidate_all`. This endpoint will 
invalidate all tokens/sessions assoiciated with the current user.

**Note:** When accessing any `/api` routes, make sure to pass the `access_token` cookie in your request. This cookie is essential for authentication and authorization purposes.

## Error Handling ❗

Here are the possible errors you might encounter while using this service:

- **INTERNAL_SERVER_ERROR:** This error indicates that something unexpected happened on the server side. If you encounter this error, please reach out to the service administrator for assistance.
- **INVALID_TOKEN:** This error occurs when the provided token is invalid, expired, or revoked. Please ensure that you are using a valid and unexpired token for your requests.
- **UNAUTHENTICATED:** This error indicates that no `access_token` cookie has been passed with the request. To access protected routes, make sure to include the `access_token` cookie containing a valid token.

Feel free to ask any questions if you need further clarification or assistance with this service. Enjoy secure and reliable authentication! 🔒✨