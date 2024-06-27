# API Documentation
# By Naufal Afthar Razzan for Coding Test on Synapsis Sinergi Digital

## Base URL: `http://localhost:9000/`

### Login

**Endpoint**: `POST /login`

**Description**: login registered users

**Example**

```json
{
    "email": "email@email.com",
    "password": "Abc123% bca"
}
```

**Reponse Success**: 200 OK

**Success Body Response**:
```json
{
    "message": "acc_token"
    // acc_token is the jwt token used Auth header on Bearer
}
```

### Register

**Endpoint**: `POST /register`

**Description**: register new users

**Example**

```json
{
    "full_name": "ujang sinaga",
    "email": "ujang_sinaga@email.com",
    "password": "Abc123% bca"
}
```

**Reponse Success**: 200 OK

**Success Body Response**:
```json
{
    "message": "new user created"
}
```

### Fetch Products Based on Category

**Endpoint**: `POST /v1/listProduct?category={category}`

**Description**: Fetch products based on its category

**Reponse Success**: 200 OK

**Success Body Response (if it has product based on category)**:
```json
{
    "message": "fetch x products based on y",
    "data": []
    // x and y are numbers of products and category respectively
    // it will returns an array of products
}
```

**Success Body Response (if it has no product based on category)**:
```json
{
    "message": "no products data of y",
    // y is category
}
```

### Insert to Shopping Cart

**Endpoint**: `POST /insertShoppingCart`

**Description**: add to shopping cart

**Example**

```json
{
    "product_id": "PROD001",
    "user_name": "ujang sinaga",
    "amount": 10
}
```

**Reponse Success**: 200 OK

**Success Body Response**:
```json
{
    "message": "success insert to shopping cart"
}
```

### Fetch Shopping Cart

**Endpoint**: `POST /v1/listShoppingCart?name={user_name}`

**Description**: Fetch a list of shopping carts based on specified user's name

**Reponse Success**: 200 OK

**Success Body Response (if a user has transaction history)**:
```json
{
    "message": "fetched x",
    "data": []
    // x is the number of transactions
    // it will returns an array of transactions
}
```

**Success Body Response (if a user has no transaction history)**:
```json
{
    "message": "you have not make any transactions",
    // y is category
}
```

### Delete Shopping Cart

**Endpoint**: `POST /v1/listShoppingCart?id={trx_id}`

**Description**: Delete shopping cart transaction based on specified transaction id

**Reponse Success**: 200 OK

**Success Body Response (if a user has transaction history)**:
```json
{
    "message": "delete success"
}
```

### Checkout Shopping Cart and Make Payment

**Endpoint**: `POST /v1/checkout?id={trx_id}`

**Description**: Move shopping cart transaction to payment by updating has_bought field

**Reponse Success**: 200 OK

**Success Body Response (if a user has transaction history)**:
```json
{
    "message": "success"
}
```