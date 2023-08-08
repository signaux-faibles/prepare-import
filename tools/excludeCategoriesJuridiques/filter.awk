
BEGIN {
    FS=","
    while (getline < INPUT_CATEGORIES)
    {
        ARRAY[$0]="oui"
    }
    close(INPUT_CATEGORIES)
}
{
    if ($28 in ARRAY) {
        print $1
    }
}