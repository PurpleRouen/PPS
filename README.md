
# PPS (Parcours Prévention Santé)
Ce microservice développé en Golang permet d'extraire les données des attestations délivrées par FFA (Fédération Française d'Athlétisme) afin de connaître les informations saisies par les coureurs et leurs validités. Afin de récupérer ces informations, il suffit d'envoyer le fichier au format PDF en tant que champs `certificate` d'un formulaire à l'endpoint `/check-pdf`.

Une fonctionnalité restante à être implémentée est la vérification de l'attestation auprès des serveurs de la FFA.

Exemple d'intégration avec PHP (Symfony) :
```php
public function postPpsCertificate(Request $request): JsonResponse
{
    $certificateFile = $request->files->get('certificate');
    $fileHandler = fopen($certificateFile->getPathname(), 'r');

    $response = $httpClient->request('POST', 'http://127.0.0.1:8080/check-pdf', [
        'headers' => [
            'Content-Type' => 'multipart/form-data',
        ],
        'body' => [
            'certificate' => $fileHandler,
        ],
    ]);
}
```

Exemple de réponse :
```json
{
    "payload": {
        "birthDate": "01/01/2000",
        "expiresAt": "08/05/2025",
        "firstName": "Romain",
        "gender": "male",
        "identifier": "P08CA013483",
        "lastName": "DUPONT"
    },
    "success": true
}
```