# E-wallet ERD

Database exercise implementing an ERD for an e-wallet.

```mermaid
---
title: E-wallet
---
erDiagram

users ||--o| users_spi : "has"

users   ||--o{ users_wallets : "has"
wallets ||--o{ users_wallets : "owned via"

transactions ||--o{ entries : "detailed by"
wallets      ||--o{ entries : "mutated by"

users {
    int id PK

    timestamp  created_at
    timestamp  updated_at
    timestamp? deleted_at

    string display_name
}

users_spi {
    int id PK

    timestamp  created_at
    timestamp  updated_at
    timestamp? deleted_at

    int        id_user     FK
    timestamp? verified_at

    string ssn        UK
    string legal_name
    date   dob

    string? tax_id
}

wallets {
    int id PK

    timestamp  created_at
    timestamp  updated_at
    timestamp? deleted_at

    bigint balance_idr
}

users_wallets {
    int id PK

    timestamp  created_at
    timestamp? deleted_at

    int id_user   FK
    int id_wallet FK
}

transactions {
    int id PK

    timestamp  created_at
    timestamp  updated_at
    timestamp? deleted_at

    enum type

    string  ref_internal
    string? ref_external
    string? provider
    string? note

    enum status
}

entries {
    int id PK

    timestamp created_at

    int id_wallet      FK
    int id_transaction FK

    bigint amount
    bigint balance_idr_after
}
```

## License

[MIT](LICENSE)
