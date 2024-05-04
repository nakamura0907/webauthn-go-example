// Base64 to ArrayBuffer
const bufferDecode = (value) => {
    return Uint8Array.from(atob(value.replace(/-/g, "+").replace(/_/g, "/")), (c) => c.charCodeAt(0));
}

// ArrayBuffer to URLBase64
const bufferEncode = (value) =>
  btoa(String.fromCharCode(...new Uint8Array(value)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");

document.getElementById("register").addEventListener("click", async function () {
    const response = await fetch("/webauthn/register/begin", {
        method: 'POST'
    });
    const data = await response.json();
    console.log(data)

    data.publicKey.challenge = bufferDecode(data.publicKey.challenge);
    data.publicKey.user.id = bufferDecode(data.publicKey.user.id);
    if (data.publicKey.excludeCredentials) {
        data.publicKey.excludeCredentials.forEach((item) => {
            item.id = bufferDecode(item.id);
        })
    }

    const credential = await navigator.credentials.create({
        publicKey: data.publicKey
    })
    console.log(credential);

    const attestationObject = credential.response.attestationObject;
    const clientDataJSON = credential.response.clientDataJSON;
    const rawId = credential.rawId;

    const finishResponse = await fetch("/webauthn/register/finish", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            id: credential.id,
            rawId: bufferEncode(rawId),
                      type: credential.type,
          response: {
            attestationObject: bufferEncode(attestationObject),
            clientDataJSON: bufferEncode(clientDataJSON),
          },
        })
    })
    const finishData = await finishResponse.json();
    console.log(finishData)
});

document.getElementById("login").addEventListener("click", async function() {
    const response = await fetch("/webauthn/login/begin", {
        method: 'POST'
    });
    const data = await response.json();
    console.log(data)

    data.publicKey.challenge = bufferDecode(data.publicKey.challenge);
    data.publicKey.allowCredentials.forEach((item) => {
        item.id = bufferDecode(item.id);
    })

    const credential = await navigator.credentials.get({
        publicKey: data.publicKey
    })
    console.log(credential);

    const assertion = credential.response;
    const clientDataJSON = assertion.clientDataJSON;
    const authenticatorData = assertion.authenticatorData;
    const signature = assertion.signature;
    const userHandle = assertion.userHandle;

    const finishResponse = await fetch("/webauthn/login/finish", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            id: credential.id,
            rawId: bufferEncode(credential.rawId),
            type: credential.type,
            response: {
                clientDataJSON: bufferEncode(clientDataJSON),
                authenticatorData: bufferEncode(authenticatorData),
                signature: bufferEncode(signature),
                userHandle: bufferEncode(userHandle)
            }
        })
    })
    const finishData = await finishResponse.json();
    console.log(finishData)
});