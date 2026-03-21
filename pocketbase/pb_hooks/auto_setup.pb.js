/// <reference path="../pb_data/types.d.ts" />

onBootstrap((e) => {
  // --- Google OAuth2 ---
  const googleClientId = $os.getenv("GOOGLE_OAUTH_CLIENT_ID")
  const googleClientSecret = $os.getenv("GOOGLE_OAUTH_CLIENT_SECRET")

  if (googleClientId && googleClientSecret) {
    try {
      const usersCollection = e.app.findCollectionByNameOrId("users")
      usersCollection.oauth2.enabled = true

      const [, found] = usersCollection.oauth2.getProviderConfig("google")
      if (found) {
        for (let i = 0; i < usersCollection.oauth2.providers.length; i++) {
          if (usersCollection.oauth2.providers[i].name === "google") {
            usersCollection.oauth2.providers[i].clientId = googleClientId
            usersCollection.oauth2.providers[i].clientSecret = googleClientSecret
            break
          }
        }
      } else {
        usersCollection.oauth2.providers.push({ name: "google", clientId: googleClientId, clientSecret: googleClientSecret })
      }

      e.app.save(usersCollection)
      console.log("[auto_setup] Google OAuth2 configured")
    } catch (err) {
      console.error("[auto_setup] Failed to configure Google OAuth2:", err)
    }
  } else {
    console.log("[auto_setup] GOOGLE_OAUTH_CLIENT_ID/SECRET not set, skipping OAuth config")
  }

  e.next()
})
