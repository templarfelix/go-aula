-- SCHEMA
create table public.categories
(
    name       text
        constraint idx_categories_name
            unique,
    id         bigint default unique_rowid() not null
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

alter table public.categories
    owner to aula;

create index idx_categories_deleted_at
    on public.categories (deleted_at);
-- FIM SCHEMA

-- BUSCAR
select * from categories;
select name, id from categories;
select name, id from categories where name = 'xxx';
select * from categories where id = 853867592319893505;
-- INSERIR
-- fixme case sensitive no name
insert into categories (name, id, created_at, updated_at, deleted_at)
values ('fritura', unique_rowid(), null, null, null);
-- ATUALIZAR
UPDATE categories SET name = 'banana'
    WHERE name = 'fritura';
-- REMOVER
DELETE FROM categories WHERE name = 'banana';

-- PRODUCT
create table public.products
(
    name       text
        constraint idx_categories_name
            unique,
    id         bigint default unique_rowid() not null
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    category bigint,
    CONSTRAINT fk_category FOREIGN KEY (category) REFERENCES categories
);

select * from products;

insert into products (name, id, created_at, updated_at, deleted_at, category)
values ('banana4', unique_rowid(), null, null,  null, 855570583266918401);

--
select * from products where category =
    (select ca.id from categories ca where ca.name = 'banana');
select * from products where category =
     (select ca.id from categories ca where ca.name = 'banana')
    AND
    deleted_at is null
    ;
select * from products where category =
                             (select ca.id from categories ca where ca.name = 'banana')
                         AND
    deleted_at is null
    and created_at >= '2023-04-10 00:00:00'
    and created_at <= '2023-04-10 23:59:59'
;