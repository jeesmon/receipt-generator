# receipt-generator

A super simple way to programmatically generate pdf donation receipts. Program takes CSV files with payment and projects to generate receipts. Sample `payments.csv` and `projects.csv` files included in the repo. Configurations for the program is read from a config yaml file (default is `config.yaml`).

# Build

```
make all
```

# Run

```
release/receipt-generator-<os>-<version> -config config.yaml
```

# Credits

* https://github.com/johnfercher/maroto
* https://socketloop.com/tutorials/golang-how-to-convert-a-number-to-words
