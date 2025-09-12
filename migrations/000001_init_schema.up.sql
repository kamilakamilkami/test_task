CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     email TEXT UNIQUE NOT NULL,
                                     password TEXT NOT NULL,
                                     role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
                                        id UUID PRIMARY KEY,
                                        name TEXT NOT NULL,
                                        price NUMERIC(10, 2) NOT NULL,
    quantity INTEGER NOT NULL
    );
