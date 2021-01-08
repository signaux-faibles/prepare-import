#!/usr/bin/python2
# -*- coding: utf-8 -*-

import json, csv
import sys, os
from datetime import datetime

def main():
  if len(sys.argv) != 2:
    print "goup.py permet de lister le contenu d'un répertoire de stockage goup"
    print "usage: goupy.py [path of directory]"
    sys.exit(1)

  if not (os.path.isdir(sys.argv[1])):
    print "Le chemin indiqué n'est pas un répertoire"
    sys.exit(1)

  data = [read_info(os.path.join(sys.argv[1], path)) for path in os.listdir(sys.argv[1]) if  read_info(os.path.join(sys.argv[1], path)) != None]

  writer = csv.DictWriter(sys.stdout, get_keys(data), dialect=csv.excel)
  writer.writeheader()
  for d in data:
    row = {k: unicode(v).encode("utf-8") for k,v in d.iteritems()}
    writer.writerow(row)

def read_info(path):
  if path[-4:] == 'info':
    with open(path) as json_file:
      try:
        data = json.load(json_file)
      except:
        return None

      data['lastModificationDate'] = datetime.fromtimestamp(os.path.getmtime(path)).isoformat()
      for k in data['MetaData'].keys():
        data['md-' + k] = data['MetaData'][k]
        del data['MetaData'][k]
      del data['MetaData']

      return data

def get_keys(datas):
  keys = set()

  for d in datas:
    keys = keys.union(d.keys())

  return keys

if __name__=="__main__":
  main()
