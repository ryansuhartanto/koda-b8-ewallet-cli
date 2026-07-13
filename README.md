# E-wallet ERD

Database exercise implementing an ERD for an e-wallet.

```mermaid
---
title: E-wallet
---
erDiagram

users {
    int id PK

    date  created_at
    date  updated_at
    date? deleted_at

    int id_spi    FK
    int id_wallet FK

    bool status_verified

}

users_spi {
    int id PK

    date  created_at
    date  updated_at
    date? deleted_at

    date verified_at

    string ssn        UK
    string legal_name
    date   dob

    string? tax_id
}

wallets {
    int id PK

    date  created_at
    date  updated_at
    date? deleted_at

    bigint balance
}

transactions {
    int id PK

    date  created_at
    date  updated_at
    date? deleted_at

    int enum_type

    string  ref_internal
    string? ref_exernal
    string? provider
    string? note
}

entries {
    int id PK

    int id_wallet      FK
    int id_transaction FK

    bigint amount
    bigint balance_after
}
```

## License

[MIT](LICENSE)
