#!/bin/bash
set -e
INPUT_SIRENE="StockUniteLegale_utf8.zip"
INPUT_FILTER="$1"

if [ -z ${1+x} ]
then
    echo "usage: exclude.sh [filter_siren]"
    exit 1
fi

if [[ ! -f ${INPUT_SIRENE} ]]
then
    echo "Le fichier StockUniteLegale_utf8.zip n'est pas disponible dans le répertoire de travail"
    exit 1
fi

if [[ ! -f ${INPUT_FILTER} ]]
then
    echo "Le fichier ${1} n'est pas disponible dans le répertoire de travail"
    exit 1
fi

unzip -p "${INPUT_SIRENE}"| awk -v INPUT_CATEGORIES="categoriesJuridiques.txt" -f filter.awk > /tmp/exclude_siren.txt
cat ${INPUT_FILTER} |awk -f exclude.awk
