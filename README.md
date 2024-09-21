# L'heure de Catherine (Catherine's hour)

I'm a simple bot for bluesky (https://bsky.app/profile/lheure2catherine.bsky.social).  
I'm just publishing a message saying "C'est l'heure de Catherine".  
I'm running on a GCP Cloud Run functions triggered by a Cloud Scheduler through a Pub/Sub topics event.  

## Configuration

I'm using the Secret Manager to store the credentials.  
It needs to be mounted as a file with this path: `/mnt/credential`  

The file is a json file with the same payload used by the `com.atproto.server.createSession` endpoint.  
It should look like this:

    {
        "identifier": "a bluesky handle",
        "password": "(app) password of the account"
    }

You should use an "app password" as this will limit the possible action I could make
(yes I can only post a fixed message).  

## Command

I have a simple command which mostly allow to locally test myself and validate that the credential are good.
