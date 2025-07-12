-- Migration: 001_create_orders_table.sql
-- Descrição: Cria a tabela orders para armazenar pedidos
-- Data: 2025-01-26

-- Criar a tabela orders
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY COMMENT 'ID único da order (UUID)',
    price DECIMAL(10,2) NOT NULL COMMENT 'Preço base da order',
    tax DECIMAL(10,2) NOT NULL COMMENT 'Taxa aplicada à order',
    final_price DECIMAL(10,2) NOT NULL COMMENT 'Preço final calculado (price + tax)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Data de criação',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Data de atualização'
);

-- Criar índice para melhorar performance das consultas
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_price ON orders(price);

-- Comentários explicativos:
-- 1. id: VARCHAR(36) para suportar UUIDs (formato: 550e8400-e29b-41d4-a716-446655440000)
-- 2. DECIMAL(10,2): Suporta valores monetários com 2 casas decimais
-- 3. Índices: Melhoram performance para ordenação e filtros
-- 4. Timestamps: Para auditoria e controle de quando os registros foram criados/modificados 