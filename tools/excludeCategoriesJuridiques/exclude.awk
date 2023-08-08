
BEGIN {
    while (getline < "/tmp/exclude_siren.txt")
        {
            ARRAY[$0]="oui"
        }
    close(INPUT_CATEGORIES)
}
{
    if ($0 in ARRAY == 0) {
        print $1
    }
}
