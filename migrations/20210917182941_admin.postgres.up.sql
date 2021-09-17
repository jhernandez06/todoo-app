CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
INSERT INTO users 
(id,first_name,last_name,email,status_user,rol,password_hash,created_at,updated_at) 
VALUES 
(uuid_generate_v4(),'Javier', 'Hernandez','jhernandez@wawand.co','activated','admin',crypt('javier', gen_salt('bf')), NOW(), NOW());