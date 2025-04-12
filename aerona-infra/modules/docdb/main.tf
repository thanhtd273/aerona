

resource "aws_security_group" "docdb_sg" {
  name        = "docdb_sg"
  description = "Allow DocumentDB traffic"
  vpc_id      = var.vpc_id

  tags = {
    Name = "docdb_sg"
  }
}

resource "aws_vpc_security_group_egress_rule" "allow_docdb_egress" {
  security_group_id = aws_security_group.docdb_sg.id
  from_port         = 0
  to_port           = 0
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_ingress_rule" "allow_docdb_ingress" {
  security_group_id = aws_security_group.docdb_sg.id
  from_port         = 27017
  to_port           = 27017
  ip_protocol       = "tcp"
  cidr_ipv4         = var.vpc_cidr
}

resource "aws_db_subnet_group" "docdb_subnet_group" {
  name       = "docdb_subnet_group"
  subnet_ids = [var.subnet_id]
  tags = {
    Name = "docdb-subnet-group"
  }
}

resource "aws_docdb_cluster_instance" "cluster_instances" {
  count              = 2
  identifier         = "docdb_cluster_${count.index}"
  cluster_identifier = aws_docdb_cluster.ticket_db.id
  instance_class     = var.db_instance_class
  engine             = "docdb"
}

resource "aws_docdb_cluster" "ticket_db" {
  cluster_identifier      = "ticket-db"
  engine                  = "docdb"
  master_username         = var.db_username
  master_password         = var.db_password
  backup_retention_period = 5
  preferred_backup_window = "07:00-09:00"
  skip_final_snapshot     = true
  vpc_security_group_ids  = [aws_security_group.docdb_sg.id]
  db_subnet_group_name    = aws_db_subnet_group.docdb_subnet_group.name
}
