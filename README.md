# Beego JWT Vergo 🚀 

Projet d’API REST en Go utilisant le framework **Beego**, un système d’authentification **JWT**, et une base de données **MySQL**.

---

## 📦 Prérequis

Avant de lancer le projet, vous devez avoir installé :
```bash
- **Go 1.20+**
- **MySQL 5.7+ ou MariaDB**
- **Git**
- (Optionnel) Beego CLI  
  
  go install github.com/beego/bee/v2@latest
# Instructions for initializing a new Git repository
```
## Étapes après avoir cloné le projet

1. **Cloner le dépôt GitHub**
```bash
git clone https://github.com/nassimfatnassi1999/beego-jwt-vergo.git
cd beego-jwt-vergo
```

2. **Installer les dépendances Go**
```bash
go mod tidy
```

3. **Configurer la base de données MySQL**
Créer une base :
```sql
CREATE DATABASE beego_jwt_vergo;
```
Configurer les accès dans :
```
conf/app.conf
```
Exemple :
```ini
db_host = localhost
db_port = 3306
db_user = root
db_password = votre_password
db_name = beego_jwt_vergo
```

4. **Générer les clés JWT si elles sont absentes**
```bash
openssl genrsa -out keys/private.txt 2048
openssl rsa -in keys/private.txt -pubout -out keys/public.txt
```

5. **Lancer le projet**
➡️ : compiler un binaire
```bash
bee run
```

Le serveur Beego se lance par défaut sur :  
👉 http://localhost:8080

