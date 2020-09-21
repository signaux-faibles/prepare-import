# Procédure pour importer les données mensuelles

## Mettre à jour les outils

Depuis `ssh stockage -R 1080` (pour se connecter à `stockage` en partageant la connexion internet de l'hôte via le port `1080`):

```sh
curl google.com --proxy socks5h://127.0.0.1:1080 # pour tester le bon fonctionnement du proxy http
git config --global http.proxy 'socks5h://127.0.0.1:1080' # si nécéssaire: pour que git utilise le proxy
cd /home/centos/prepare-import
git pull # pour mettre à jour les outils
go build
```

## Mettre les nouveaux fichiers dans un répertoire spécifique

Depuis `ssh stockage`:

```sh
sudo su
cd /var/lib/goup_base/public
mkdir _<batch>_
# /!\ Attention commande suivante non fonctionnelle !
find -maxdepth 1 -ctime -10 -print0 | xargs -0 mv -t _<batch>_/
```

## Nettoyer les fichiers de l'urssaf

Convertir les fichiers en encodage UTF-8 avant l'usage de `filter_unprintable.sh`:

```sh
iconv -f ISO_8859-1 -t UTF-8
```

Afficher puis supprimer les lignes illisibles, depuis `ssh stockage`:

```sh
cd _<batch>_
# Print
sh /home/centos/goup_transfert/public/tools/filter_unprintable.sh -p *
# Replace
sh /home/centos/goup_transfert/public/tools/filter_unprintable.sh -r *
```

## Créer le filtre basé sur le fichier effectif

Afin d'importer les données, il faut commencer par créer le périmètre des
entreprises intéressantes à partir du fichier effectif de l'ACOSS.

Depuis `ssh stockage`:

```sh
cd `batch`
create_filter --path ./effectif.csv
```

## Télécharger le fichier Siren

Depuis `ssh stockage -R 1080` (avec partage de connexion internet de l'hôte via le port `1080`):

```sh
http_proxy="socks5h://127.0.0.1:1080" wget http://data.cquest.org/geo_sirene/v2019/last/StockEtablissement_utf8_geo.csv.gz
https_proxy="socks5h://127.0.0.1:1080" wget https://www.data.gouv.fr/fr/datasets/r/c63c91ec-7659-490b-baac-98ee599ece37
```

Note: penser à mettre les URLs à jour.

## Télécharger le fichier Diane

1. Se connecter sur le site [Diane+](https://diane.bvdinfo.com)

2. _Créer un fichier de filtrage à partir du fichier effectif._
Regarder le numéro de la nouvelle variable à importer (le suivant du dernier
numéro déjà importé dans:
_Mes données_ > _Données importées_ > _Importer nouvelle variable_

3. Changer le fichier `filter_to_diane.awk`
pour mettre à jour le numéro de variable.
Par exemple si le dernier est CF00011 dans diane+ alors il faut mettre CF00012
dans le script.
/!\ Attention, le script n'est pas robuste, par exemple si la sélection de
département est décommentée, il faut changer l'encodage et le séparateur de la
commande suivante non commentée (options -e et -d) /!\

4. Créer la nouvelle variable en indiquant qu'il s'agit d'un champs `identifiant d'entreprise`
Récupérer le fichier sur l'ordinateur local, le transformer en fichier excel,
et le soumettre sur diane+ dans l'interface _importer nouvelle variable_

5. Sélectionner la nouvelle variable dans:
_Mes données_ > _Données importées_ > _Entreprises avec une donnée importée_

> _Autres ..._

Cette sélection peut-être complétée avec:
_Entreprises mises à jour_ > _Données financières et descriptives_

## Créer un objet admin pour l'intégration des données

Utiliser `prepare-import` depuis `ssh stockage`:

```sh
~/prepare-import/prepare-import -batch "<BATCH>" -date-fin-effectif "<DATE>" -path "../goup/public"
```

- Ne pas oublier le fichier filter !
- Il faut également aller consulter à la main la dernière colonne non vide du
  fichier effectif et renseigner sa valeur dans le fichier admin.

- Et enfin changer le nom du batch en langage naturel: ex "Février 2020".

## Lancer l'import

Depuis `ssh stockage -t tmux att`:

```sh
export http_proxy=""; http :3000/api/data/import batch="2002_1"
```

Vérifier les logs que l'import s'est bien passé.

## Lancer le compactage

Depuis `ssh stockage -t tmux att`:

```sh
export http_proxy=""; http :3000/api/data/compact batch="2002_1"
```