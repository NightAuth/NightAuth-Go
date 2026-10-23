# NightAuth Go SDK — Donnez-nous une étoile sur GitHub ⭐

SDK officiel **NightAuth** pour Go et suite de clients de référence (interface graphique moderne en **Win32 GDI natif** et CLI interactive pour terminal). Protégez vos logiciels Go contre le piratage, l'ingénierie inverse et la falsification avec gestion des licences cloud, verrouillage HWID, chat chiffré en temps réel et téléchargement sécurisé en mémoire en quelques minutes.

Portail Officiel : [https://nightauth.org](https://nightauth.org)

---

**Languages / Idiomas / Langues / Sprachen:**  
[English](../README.md) • [Español](es.md) • [Français](fr.md) • [Deutsch](de.md) • [Italiano](it.md) • [Português](pt.md) • [Русский](ru.md) • [Türkçe](tr.md) • [Polski](pl.md) • [中文](zh.md)

---

> [!NOTE]
> ### 📺 Tutoriels Vidéo et Guides de Configuration
> Les tutoriels vidéo officiels et les guides pas à pas expliquant comment configurer vos applications, plans, abonnements, variables cloud, webhooks et licences dans le tableau de bord NightAuth **seront publiés prochainement** sur notre portail officiel et chaîne YouTube.

---

## 📸 Captures d'Écran

### Application de Bureau (Thème Sombre Moderne en Win32 GDI Natif)

| Écran d'Accueil / Authentification | Tableau de Bord & Gestion des Licences |
| :---: | :---: |
| ![NightAuth Desktop Login](../image/desktop.png) | ![NightAuth Desktop Dashboard](../image/desktop-dashboard.png) |

### Application Console (CLI Interactive pour Terminal)

| Écran de Connexion Terminal | Tableau de Bord Interactif Console |
| :---: | :---: |
| ![NightAuth Console Login](../image/console.png) | ![NightAuth Console Dashboard](../image/console-dashboard.png) |

---

## 🛡️ Pratiques de Sécurité Recommandées

* **Obfuscation Avancée en Go :** Les binaires Go standard conservent les noms de symboles et métadonnées de réflexion. Pour vos versions de production, compilez avec **Garble** (`garble -literals -tiny build`) afin de supprimer les types et chiffrer les chaînes.
* **Empaqueteurs Commerciaux :** Protégez vos exécutables avec des solutions comme **VMProtect** ou **Themida** pour vous prémunir contre IDA Pro, Ghidra et x64dbg.
* **Intégrité Mémoire & Anti-Débogage :** NightAuth intègre des protections contre les débogueurs (`x64dbg`, `CheatEngine`, `Wireshark`, `HTTPDebugger`) et proxies système. Invoquez `nightauth.CheckDebugger()` et `nightauth.CheckBlacklistedProcesses()` au démarrage (voir [SecurityProtections](#24-securityprotections)).
* **Exécution Sécurisée en Mémoire :** Lors du téléchargement de modules protégés avec `auth.DownloadFile()`, **n'écrivez jamais les octets sur disque**. Exécutez-les directement en RAM et nettoyez les mémoires tampons immédiatement après.

---

## 📦 Prérequis & Compilation

Nécessite **Go 1.20** ou supérieur.

L'application de bureau Win32 native utilise les DLL système de Windows (`user32.dll`, `gdi32.dll`, `kernel32.dll`, `msimg32.dll`) sans aucune dépendance CGO.

### Compiler le Client Console
```bash
cd console
go build -ldflags="-s -w" -o bin/NightAuth-Console.exe .
```

### Compiler le Client de Bureau (Sans Fenêtre Console)
```bash
cd desktop
go build -ldflags="-s -w -H=windowsgui" -o bin/NightAuth-Desktop.exe .
```

---

## ⚙️ Configuration & Initialisation

Configurez les identifiants de votre application depuis la section **Manage Apps** du tableau de bord NightAuth :

```go
package main

import (
	"fmt"
	"os"

	"nightauth"
)

// Configuration de l'Application
const (
	AppName    = ""
	AppApiKey  = ""
	AppVersion = "1.0"
	AppToken   = ""
)

func main() {
	// Initialiser le client API
	auth := nightauth.NewApi(AppName, AppApiKey, AppVersion, AppToken, "")

	// Établir le handshake cryptographique (Curve25519 + Ed25519 + AES-256-GCM)
	auth.Init()

	if !auth.Response.Success {
		fmt.Printf("Échec de l'initialisation : %s\n", auth.Response.Message)
		os.Exit(1)
	}

	fmt.Println("Connecté avec succès au cloud NightAuth !")
}
```

---

## 🚀 Guide d'Utilisation du SDK

### 1. Validation de Session & Vérification de Liste Noire

```go
// Vérifier la validité de la session active
if !auth.Check() {
	fmt.Printf("Session invalide : %s\n", auth.Response.Message)
}

// Vérifier si la machine ou l'adresse IP est sur liste noire
hwid := nightauth.SecurityProtections{}.GetHardwareId()
if auth.CheckBlacklist(hwid) {
	fmt.Println("Accès Refusé : Machine ou IP restreinte.")
	os.Exit(1)
}
```

### 2. Connexion avec Nom d'Utilisateur et Mot de Passe (et 2FA)

```go
if auth.Login(username, password, totpCode) {
	fmt.Printf("Bienvenue, %s !\n", auth.User.Username)
} else {
	fmt.Printf("Connexion échouée : %s\n", auth.Response.Message)
}
```

### 3. Inscription avec Clé de Licence

```go
if auth.Register(username, password, licenseKey, email) {
	fmt.Printf("Inscription réussie : %s\n", auth.User.Username)
} else {
	fmt.Printf("Échec de l'inscription : %s\n", auth.Response.Message)
}
```

### 4. Authentification par Clé de Licence Uniquement (et 2FA)

```go
// Connexion par clé de licence uniquement
if auth.License(licenseKey) {
	fmt.Printf("Licence valide ! Connecté en tant que : %s\n", auth.User.Username)
} else {
	fmt.Printf("Échec d'authentification par licence : %s\n", auth.Response.Message)
}

// Connexion par clé de licence avec code 2FA TOTP
if auth.LicenseWithTotp(licenseKey, "123456") {
	fmt.Printf("Licence vérifiée avec TOTP ! Connecté en tant que : %s\n", auth.User.Username)
}
```

### 5. Mise à Niveau d'Abonnement / Clé Supplémentaire

```go
if auth.Upgrade(username, upgradeKey) {
	fmt.Println("Abonnement mis à niveau avec succès !")
} else {
	fmt.Printf("Échec de la mise à niveau : %s\n", auth.Response.Message)
}
```

### 6. Réinitialisation de l'Empreinte Matérielle (Verrou HWID)

```go
if auth.ResetHwid(username, password, totpCode) {
	fmt.Println("HWID réinitialisé. Vous pouvez maintenant vous connecter sur cet appareil.")
} else {
	fmt.Printf("Échec de la réinitialisation HWID : %s\n", auth.Response.Message)
}
```

### 7. Récupération de Compte / Mot de Passe Oublié

```go
// Étape 1 : Demande de lien de réinitialisation
if auth.ForgotPassword(email) {
	fmt.Println("Lien de réinitialisation envoyé par e-mail.")
}

// Étape 2 : Confirmation avec jeton et nouveau mot de passe
if auth.ResetPassword(resetToken, newPassword) {
	fmt.Println("Mot de passe mis à jour avec succès !")
}
```

### 8. Données du Profil Utilisateur (`auth.User`)

```go
user := auth.User
fmt.Printf("Utilisateur : %s\n", user.Username)
fmt.Printf("E-mail :      %s\n", user.Email)
fmt.Printf("HWID :        %s\n", user.Hwid)
fmt.Printf("IP Client :   %s\n", user.Ip)
fmt.Printf("Création :    %s\n", user.CreatedAt)
fmt.Printf("Dernier log : %s\n", user.LastLoginAt)

for _, sub := range user.Subscriptions {
	fmt.Printf("Plan : %s | Expiration : %s | Clé : %s\n", sub.Subscription, sub.Expiry, sub.Key)
}
```

### 9. Variables Cloud (Globales)

```go
motd := auth.GetGlobalVar("server_motd")
fmt.Printf("Variable Cloud (Globale) : %s\n", motd)
```

### 10. Variables Personnalisées Utilisateur

```go
// Enregistrer une variable utilisateur
auth.SetUserVar("custom_theme", "dark_neon")

// Récupérer la variable enregistrée
val := auth.GetUserVar("custom_theme")
fmt.Printf("Variable Utilisateur : %s\n", val)

// Lister toutes les variables
if auth.GetUserVariables() {
	rawVars, _ := auth.Response.Payload["variables"].([]interface{})
	for _, item := range rawVars {
		if m, ok := item.(map[string]interface{}); ok {
			fmt.Printf("Variable %s = %s\n", m["key"], m["value"])
		}
	}
}
```

### 11. Journaux d'Audit Télémesure

```go
auth.Log("Tentative de violation d'intégrité interceptée.")
```

### 12. Blocage Utilisateur et Machine (Killswitch / Ban)

```go
auth.Ban("Échec d'intégrité Anti-Cheat : Modification de mémoire détectée.")
```

### 13. Webhooks Côté Serveur & Envoi de Fichiers

```go
// Webhook JSON classique
payload := map[string]interface{}{
	"event": "alert_triggered",
	"level": "critical",
}
result := auth.TriggerWebhook("7kR0UedlVI", payload)
if result != nil && result.Success {
	fmt.Printf("Webhook délivré (HTTP %d)\n", result.StatusCode)
}

// Webhook multipart avec fichier joint
fileBytes := []byte("données de rapport")
fileResult := auth.TriggerWebhookFile("7kR0UedlVI", fileBytes, "report.txt", payload)
if fileResult != nil && fileResult.Success {
	fmt.Printf("Webhook avec fichier envoyé (HTTP %d)\n", fileResult.StatusCode)
}
```

### 14. Téléchargement Sécurisé en Mémoire Vive

```go
fileBytes := auth.DownloadFile("385624")
if len(fileBytes) > 0 {
	fmt.Printf("Téléchargé %d octets en mémoire vive de manière sécurisée.\n", len(fileBytes))
	// Exécuter directement en mémoire sans toucher au disque
}
```

### 15. Salons de Discussion en Temps Réel

```go
// Envoyer un message
auth.SendChatMessage("Bonjour la communauté depuis Go !", "general")

// Récupérer les messages récents
messages := auth.GetChatMessages("general")
for _, msg := range messages {
	fmt.Printf("[%s] %s : %s\n", msg.Timestamp, msg.Author, msg.Message)
}

// Récupérer l'historique paginé
paged := auth.GetChatMessagesPage("general", 0, true)
for _, msg := range paged {
	fmt.Printf("#%s [%s] %s : %s\n", msg.Id, msg.Timestamp, msg.Author, msg.Message)
}
```

### 16. Flux d'Événements en Temps Réel (Server-Sent Events)

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err := auth.SubscribeEvents(ctx, "", func(evt nightauth.EventFrame) {
	fmt.Printf("[ÉVÉNEMENT #%s] %s\n", evt.ID, evt.Data)
})
if err != nil {
	fmt.Printf("Flux d'événements terminé : %v\n", err)
}
```

### 17. Authentification à Deux Facteurs (2FA)

```go
// Configuration 2FA (TOTP)
setupData := auth.Setup2FA()
if setupData != nil {
	fmt.Printf("Secret : %s | URI : %s\n", setupData.Secret, setupData.Uri)
}

// Confirmation d'activation
confirmData := auth.Confirm2FA("123456")
if confirmData != nil && confirmData.Success {
	fmt.Println("2FA Activé avec succès ! Codes de secours :")
	for _, code := range confirmData.BackupCodes {
		fmt.Printf(" - %s\n", code)
	}
}

// Désactivation du 2FA
if auth.Disable2FA(password, "123456") {
	fmt.Println("2FA désactivé avec succès.")
}
```

### 18. Cycle de Vie de Session & Déconnexion

```go
// Latence périodique pour maintenir la session active
auth.Heartbeat()

// Révocation de la session et fermeture
auth.Logout()
auth.Close()
```

### 19. Changer de Nom d'Utilisateur

```go
if auth.ChangeUsername("nouvel_alias", "mot_de_passe_actuel") {
	fmt.Printf("Nom d'utilisateur mis à jour avec succès : %s\n", auth.User.Username)
} else {
	fmt.Printf("Échec de modification du nom d'utilisateur : %s\n", auth.Response.Message)
}
```

### 20. Statistiques de l'Application & Utilisateurs en Ligne

```go
// Récupérer les statistiques de l'application
stats := auth.GetStats()
if stats != nil {
	fmt.Printf("Utilisateurs : %d au total, %d actifs, %d bannis\n", stats.TotalUsers, stats.ActiveUsers, stats.BannedUsers)
	fmt.Printf("Licences : %d au total, %d utilisées, %d inutilisées\n", stats.TotalLicenses, stats.UsedLicenses, stats.UnusedLicenses)
	fmt.Printf("Version serveur : %s\n", stats.Version)
}

// Récupérer la liste des utilisateurs en ligne
onlineUsers := auth.GetOnlineUsers()
fmt.Printf("Utilisateurs en ligne (%d) : %v\n", len(onlineUsers), onlineUsers)
```

### 21. Vérifier le Hash du Fichier

```go
// Vérifier le hash SHA-256 du binaire par rapport aux hashs approuvés
isValid := auth.VerifyFileHash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
if isValid {
	fmt.Println("Intégrité du fichier vérifiée avec succès.")
} else {
	fmt.Println("Échec du contrôle d'intégrité : hash non conforme.")
}
```

### 22. Tickets d'Accès

```go
// Demander un ticket d'accès temporaire avec portée
ticket := auth.AccessTicket("file_download", "file_abc123")
if ticket != "" {
	fmt.Printf("Ticket d'accès accordé : %s\n", ticket)
} else {
	fmt.Printf("Échec de création du ticket d'accès : %s\n", auth.Response.Message)
}
```

### 23. Callbacks Réactifs & Écoute en Arrière-Plan

Configurez des gestionnaires d'événements réactifs et démarrez le dispatcher en arrière-plan :

```go
auth.OnSessionInvalidated = func(reason string) {
	fmt.Printf("Session fermée : %s\n", reason)
	os.Exit(0)
}
auth.OnAppBanned = func(reason string) {
	fmt.Printf("Application bannie : %s\n", reason)
}
auth.OnAppStatusChanged = func(status string) {
	fmt.Printf("Statut de l'application mis à jour : %s\n", status)
}
auth.OnVersionUpdated = func(version string) {
	fmt.Printf("Nouvelle version disponible : %s\n", version)
}
auth.OnChatMessageReceived = func(msg nightauth.ChatMessage, channel string) {
	fmt.Printf("[%s] %s: %s\n", channel, msg.Author, msg.Message)
}

// Démarrer l'écouteur SSE en arrière-plan
auth.StartRealtimeEventListener()

// Arrêter l'écouteur lorsqu'il n'est plus nécessaire
auth.StopRealtimeEventListener()
```

### 24. `SecurityProtections`

Utilitaires statiques de sécurité et anti-altération :

```go
// Exécuter les contrôles automatiques (détection de proxy HTTP système)
nightauth.RunSecurityChecks()

// Détecter les débogueurs actifs (Win32 locaux ou distants)
nightauth.CheckDebugger()

// Détecter les outils d'analyse en liste noire (x64dbg, Cheat Engine, IDA, Wireshark, Fiddler)
// Définir NIGHTAUTH_IGNORE_SECURITY=1 lors des tests locaux
nightauth.CheckBlacklistedProcesses()

// Obtenir l'identifiant matériel stable de l'appareil
hwid := nightauth.GetHardwareId()
fmt.Printf("HWID de l'appareil : %s\n", hwid)
```

### 25. Options Avancées & Clés Serveur Personnalisées

```go
auth := nightauth.NewApiWithOptions(
	"my-app",
	"my-api-key",
	"1.0",
	"",
	"https://nightauth.org/api/v1",
	false,
)

auth.SetRequestTimeout(15 * time.Second)

err := auth.SetServerKeys(
	"cdb2285db481666ce40003aadb874c8751dbeb9a0672a37f740dfe881ecdc742",
	"b773fc4f6cf1a0938e78e714c3d84cb386584bea24728923008471a6cd203213",
)
```

### 26. Mises à Jour Signées Cryptographiquement

Vérifie et télécharge des versions inviolables signées avec Ed25519 :

```go
// Téléchargement vérifié en un coup avec protection anti-rollback
targetPath, err := nightauth.SignedUpdateInstaller.Download(
	"https://releases.example.com/manifest.json",
	"VOTRE_CLE_PUBLIQUE_PINNED_ED25519_HEX",
	"my-app",
	"updates/my-app.exe",
	5*time.Minute,
)
if err != nil {
	fmt.Printf("Échec de la mise à jour : %v\n", err)
} else {
	fmt.Printf("Mise à jour vérifiée installée dans : %s\n", targetPath)
}

// Ou pilotez manuellement le processus :
verified, err := nightauth.VerifySignedUpdate(manifestEnvelopeJSON, pinnedPublicKeyHex, "my-app", minimumSequence)
if err == nil {
	fmt.Printf("Version vérifiée %s v%s (séquence %d)\n", verified.Application, verified.Version, verified.Sequence)
	installed, err := nightauth.DownloadSignedUpdate(verified, "updates/my-app.exe", 5*time.Minute)
}
```

## 🖥️ Exécuter les Clients Console et Bureau

Lancez l'un des clients de référence depuis le répertoire `source` :

```bash
# Interactive Terminal Console Client
cd console
go run .

# Modern Pure Win32 Desktop Client (Windows)
cd desktop
go run .
```

---

## 📄 Licence

Ce projet est sous licence **Elastic License 2.0**.

* Vous ne pouvez pas fournir le logiciel à des tiers en tant que service hébergé ou géré.
* Vous ne pouvez pas déplacer, modifier, désactiver ou contourner la fonctionnalité de clé de licence.
* Vous ne pouvez pas altérer, supprimer ou masquer les mentions de licence ou de droit d'auteur.

Consultez le fichier [LICENSE](../LICENSE) pour le texte juridique complet.
