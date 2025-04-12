resource "aws_security_group" "postgres_sg" {
  name        = "postgres-sg"
  description = "Allow PostgreSQL traffic"
  vpc_id      = var.vpc_id

  tags = {
    Name = "postgres-sg"
  }
}

resource "aws_vpc_security_group_egress_rule" "allow_postgres_egress" {
  security_group_id = aws_security_group.postgres_sg.id
  from_port         = 0
  to_port           = 0
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_ingress_rule" "allow_postgres_ingress" {
  security_group_id = aws_security_group.postgres_sg.id
  from_port         = 5432
  to_port           = 5432
  ip_protocol       = "tcp"
  cidr_ipv4         = var.vpc_cidr
}

resource "aws_db_subnet_group" "postgres_subnet_group" {
  name       = "postgres-subnet-group"
  subnet_ids = [var.subnet_id]
  tags = {
    Name = "postgres-subnet-group"
  }
}

resource "aws_db_instance" "booking_db" {
  identifier             = "booking-db"
  allocated_storage      = 10
  db_name                = "bookingdb"
  engine                 = "postgres"
  engine_version         = var.db_engine_version
  instance_class         = var.db_instance_class
  username               = var.db_username
  password               = var.db_password
  vpc_security_group_ids = [aws_security_group.postgres_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.postgres_subnet_group.name
  skip_final_snapshot    = true
  publicly_accessible    = false
  timeouts {
    create = "3h"
    delete = "3h"
    update = "3h"
  }

  tags = {
    Name = "booking-db"
  }
}

resource "aws_db_instance" "payment_db" {
  identifier             = "payment-db"
  allocated_storage      = 10
  db_name                = "paymentdb"
  engine                 = "postgres"
  engine_version         = var.db_engine_version
  instance_class         = var.db_instance_class
  username               = var.db_username
  password               = var.db_password
  vpc_security_group_ids = [aws_security_group.postgres_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.postgres_subnet_group.name
  skip_final_snapshot    = true
  publicly_accessible    = false
  timeouts {
    create = "3h"
    delete = "3h"
    update = "3h"
  }

  tags = {
    Name = "payment-db"
  }
}
