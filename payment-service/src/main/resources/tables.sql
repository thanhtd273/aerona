create table payment (
     payment_id varchar(50) primary key,
     booking_id varchar(50) not null,
     pnr varchar(15) not null,
     payment_method varchar(15) not null,
     external_payment_id varchar(50),
     amount float not null,
     currency varchar(10) not null,
     status varchar(15) not null,
     created_at timestamp default current_timestamp,
     paid_at timestamp,
     failed_at timestamp,
     refunded_at timestamp,
     failure_reason varchar(255)
);

create index idx_booking_id on payment(booking_id);
create index idx_pnr on payment(pnr);
create index idx_external_payment_id on payment(external_payment_id);

create table payment_log (
     id varchar(50) primary key,
     payment_id varchar(50) not null,
     status varchar(15) not null,
     event_time timestamp not null default current_timestamp,
     description varchar(255),
     constraint payment_fk foreign key (payment_id) references payment(payment_id)
);
create index idx_payment_id on payment_log(payment_id);

create table refund (
    refund_id varchar(50) primary key,
    payment_id varchar(50) not null,
    external_refund_id varchar(50),
    amount float not null,
    currency varchar(10) not null,
    requested_at timestamp not null default current_timestamp,
    completed_at timestamp,
    status varchar(15) not null,
    reason varchar(255),
    constraint payment_fk foreign key (payment_id) references payment(payment_id)
);
create index idx_payment_refund_id on refund(payment_id);
create index inx_external_refund_id on refund(external_refund_id);