-- Migration: 003_seed_data.sql
-- Dados iniciais para desenvolvimento

-- Inserir usuários de exemplo
INSERT INTO users (id, full_name, document, email, password, user_type, balance) VALUES
(
    '550e8400-e29b-41d4-a716-446655440001',
    'João Silva',
    '11144477735',
    'joao@teste.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: 123456
    'common',
    1000.00
),
(
    '550e8400-e29b-41d4-a716-446655440002',
    'Maria Santos',
    '22244477735',
    'maria@teste.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: 123456
    'common',
    500.00
),
(
    '550e8400-e29b-41d4-a716-446655440003',
    'Loja do João LTDA',
    '11222333000181',
    'loja@teste.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: 123456
    'merchant',
    0.00
),
(
    '550e8400-e29b-41d4-a716-446655440004',
    'Pedro Oliveira',
    '33344477735',
    'pedro@teste.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: 123456
    'common',
    2500.00
),
(
    '550e8400-e29b-41d4-a716-446655440005',
    'Supermercado ABC LTDA',
    '22333444000195',
    'supermercado@teste.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', -- password: 123456
    'merchant',
    0.00
);

-- Inserir algumas transações de exemplo
INSERT INTO transactions (id, payer_id, payee_id, amount, status, authorization_id, notification_sent, completed_at) VALUES
(
    '660e8400-e29b-41d4-a716-446655440001',
    '550e8400-e29b-41d4-a716-446655440001', -- João
    '550e8400-e29b-41d4-a716-446655440003', -- Loja do João
    100.00,
    'completed',
    'auth_123456789',
    true,
    CURRENT_TIMESTAMP - INTERVAL '2 days'
),
(
    '660e8400-e29b-41d4-a716-446655440002',
    '550e8400-e29b-41d4-a716-446655440004', -- Pedro
    '550e8400-e29b-41d4-a716-446655440002', -- Maria
    250.00,
    'completed',
    'auth_987654321',
    true,
    CURRENT_TIMESTAMP - INTERVAL '1 day'
),
(
    '660e8400-e29b-41d4-a716-446655440003',
    '550e8400-e29b-41d4-a716-446655440001', -- João
    '550e8400-e29b-41d4-a716-446655440005', -- Supermercado
    50.00,
    'pending',
    null,
    false,
    null
);
