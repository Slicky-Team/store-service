CREATE TYPE user_role AS ENUM ('ADMIN', 'MANAGER', 'STAFF');

CREATE TABLE user_profiles (
    id UUID PRIMARY KEY,
    account_id UUID UNIQUE NOT NULL,
    full_name TEXT NOT NULL,
    age INTEGER,
    phone_number TEXT,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Table: brands
CREATE TABLE brands (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    metadata JSON,
    rating FLOAT DEFAULT 0,
    owner_id UUID NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_brand_owner FOREIGN KEY (owner_id) REFERENCES user_profiles(id) ON DELETE CASCADE
);
CREATE INDEX idx_brands_owner_id ON brands(owner_id);
CREATE INDEX idx_brands_is_active ON brands(is_active);

-- Table: stores
CREATE TABLE stores (
    id UUID PRIMARY KEY,
    brand_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    metadata JSON,
    rating FLOAT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_store_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
);
CREATE INDEX idx_stores_brand_id ON stores(brand_id);
CREATE INDEX idx_stores_is_active ON stores(is_active);

-- Table: store_services
CREATE TABLE store_services (
    id UUID PRIMARY KEY,
    store_id UUID NOT NULL,
    service_name TEXT NOT NULL,
    service_price FLOAT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    metadata JSON,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_service_store FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE
);
CREATE INDEX idx_services_store_id ON store_services(store_id);
CREATE INDEX idx_services_is_active ON store_services(is_active);

-- Table: user_store
CREATE TABLE user_store (
    id UUID PRIMARY KEY,
    store_id UUID NOT NULL,
    user_id UUID NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user_store FOREIGN KEY (store_id) REFERENCES stores(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_store_user FOREIGN KEY (user_id) REFERENCES user_profiles(account_id) ON DELETE CASCADE
);
CREATE INDEX idx_user_store_store_id ON user_store(store_id);
CREATE INDEX idx_user_store_user_id ON user_store(user_id);

-- Table: discounts
CREATE TABLE discounts (
    id UUID PRIMARY KEY,
    brand_id UUID NOT NULL,
    discount_name TEXT NOT NULL,
    percentage FLOAT CHECK (percentage >= 0 AND percentage <= 100),
    metadata JSON,
    is_active BOOLEAN DEFAULT TRUE,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_discount_brand FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE
);
CREATE INDEX idx_discounts_brand_id ON discounts(brand_id);
CREATE INDEX idx_discounts_is_active ON discounts(is_active);
