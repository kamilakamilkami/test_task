DELETE FROM departments
WHERE id = 'f0588c10-2160-4f73-b2fe-1d65519d8285';

DELETE FROM document_types
WHERE id = '2df98bb8-de8f-4ea1-836f-149ae6999a05';

DELETE FROM employees WHERE id IN ('11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', '33333333-3333-3333-3333-333333333333');

DELETE FROM users WHERE email IN ('kamilafake@gmail.com', 'rizafake@gmail.com', 'alishfake@gmail.com');




INSERT INTO departments (id, name, code, parent_id)
VALUES ('f0588c10-2160-4f73-b2fe-1d65519d8285', 'ИТ отдел', 'IT', NULL)
ON CONFLICT (id) DO NOTHING;


DELETE FROM document_types WHERE code='CERT';

INSERT INTO document_types (id, code, name, workflow)
VALUES ('2df98bb8-de8f-4ea1-836f-149ae6999a00', 'CERT', 'Справка с места работы', '')
ON CONFLICT (id) DO NOTHING;


DELETE FROM employees
WHERE id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333'
);

INSERT INTO employees (
    id, fio, iin, email, phone, birth_date, employed_at, terminated_at, status,
    department_id, position_id, grade, employment_type, salary_base, salary_currency,
    work_schedule, manager_id, created_at
) VALUES
('11111111-1111-1111-1111-111111111111', 'Камиля', '111111111111', 'kamilafake@gmail.com', '77000000001',
 '1995-05-05', '2022-01-01', NULL, 'ACTIVE',
 'f0588c10-2160-4f73-b2fe-1d65519d8285', '00000000-0000-0000-0000-000000000000',
 'A', 'FULL', 120000.00, 'RUB', '5/2', NULL, NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO employees (id, fio, iin, email, phone, birth_date, employed_at, status,
    department_id, position_id, grade, employment_type, salary_base, salary_currency,
    work_schedule, created_at)
VALUES
('22222222-2222-2222-2222-222222222222', 'Риза', '222222222222', 'rizafake@gmail.com', '77000000002',
 '1996-06-06', '2022-02-01', 'ACTIVE',
 'f0588c10-2160-4f73-b2fe-1d65519d8285', '00000000-0000-0000-0000-000000000000',
 'B', 'FULL', 130000.00, 'RUB', '5/2', NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO employees (id, fio, iin, email, phone, birth_date, employed_at, status,
    department_id, position_id, grade, employment_type, salary_base, salary_currency,
    work_schedule, created_at)
VALUES
('33333333-3333-3333-3333-333333333333', 'Алиш', '333333333333', 'alishfake@gmail.com', '77000000003',
 '1997-07-07', '2022-03-01', 'ACTIVE',
 'f0588c10-2160-4f73-b2fe-1d65519d8285', '00000000-0000-0000-0000-000000000000',
 'C', 'FULL', 140000.00, 'RUB', '5/2', NOW())
ON CONFLICT (id) DO NOTHING;

DELETE FROM users
WHERE id IN (
    'e567d8f6-e816-4c87-b209-d1d4ea4f9594',
    'bd59e3d8-597d-4c95-ab92-124e991fe8e8',
    '11786bd2-ff07-4478-a927-0fb5fafbb39f'
);


INSERT INTO users (id, name, email, password, role, employee_id)
VALUES
('e567d8f6-e816-4c87-b209-d1d4ea4f9594', '', 'kamilafake@gmail.com', '$2a$10$2eT4Va0C/o1ptB1MHhljlO3CM0QCJajMg0FL2xM.6TbEW3WbjTuW2', 'EMPLOYEE', '11111111-1111-1111-1111-111111111111'),
('bd59e3d8-597d-4c95-ab92-124e991fe8e8', '', 'rizafake@gmail.com', '$2a$10$/bvQZDWh2rnVa/Tx3rgDQ.nqi0zaTjnS6Ug9zO8g5KBzS3j52m.V6', 'ADMIN', '22222222-2222-2222-2222-222222222222'),
('11786bd2-ff07-4478-a927-0fb5fafbb39f', '', 'alishfake@gmail.com', '$2a$10$CCtl1iu9rT5..hGI4GMTDeVNTLLeK2LZnL0Tsh2i..bWp9hnoZX8e', 'MANAGER', '33333333-3333-3333-3333-333333333333')
ON CONFLICT (id) DO NOTHING;

DELETE FROM templates
WHERE id IN (
    '46e55e17-5d97-46ab-b2b3-2cf41f2d5e9d'
);

INSERT INTO templates (
    id, name, type, engine, body, placeholders, version, is_active, created_by, created_at
) VALUES (
    '46e55e17-5d97-46ab-b2b3-2cf41f2d5e9d',
    'Справка с места работы',
    'CERTIFICATE',
    'HTML',
'<html>
<head>
<meta charset="UTF-8">
<title>Справка</title>
</head>
<body>
<p>Настоящая справка выдана {{employee.fullName}}, работающему(ей) в {{company_name}}</p>
<p>В должности {{employee.position}} с {{employee.hireDate}}.</p>
{{#if includeSalary}}
<p>Оклад: {{employee.salaryBase}} {{employee.salaryCurrency}} ({{salaryInWords}}).</p>
{{/if}}
<p>Дата: {{date}}</p>
</body>
</html>',
    '["employee.fullName", "employee.position", "employee.hireDate", "company.name", "certificate.includeSalary", "employee.salaryBase", "employee.salaryCurrency", "salaryInWords", "date.format"]',
    2,
    true,
    '33333333-3333-3333-3333-333333333333',
    '2025-09-13 02:51:20.009967'
)
ON CONFLICT (id) DO NOTHING;

