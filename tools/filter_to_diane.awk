#!/usr/bin/awk -f

# Usage: Transforms a list of sirets into a diane format csv file. 
# 
# ./filter_to_diane -v var_num="CF000xx" ../20xx/filter_20xx.csv > ../diane_req/diane_filter_20xx.csv
# 
# 
# It then needs to be transfered to local: 
# scp stockage:/var/lib/goup_base/public/diane_req/diane_filter_20xx.csv /home/pierre/go/src/github.com/signaux-faibles/opensignauxfaibles/data-raw/diane/req_diane
#  
# And converted to xlsx BY HAND as ssconvert transforms siren text into numbers and drops heading 0s. 
# 


BEGIN { OFS=","; print var_num, "SIREN" }
{if (!siren[$0]) { i++; print "id"i, "\x27"$0; siren[$0]=1 }}
