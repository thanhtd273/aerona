resource "aws_security_group" "allow_ssh" {
  name        = "Allow SSH"
  description = "Allow SSH to developers"
  vpc_id      = var.vpc_id

  tags = {
    Name = "Allow SSH"
  }
}

resource "aws_vpc_security_group_ingress_rule" "allow_ssh_ipv4" {
  security_group_id = aws_security_group.allow_ssh.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 22
  ip_protocol       = "tcp"
  to_port           = 22
}

resource "aws_vpc_security_group_egress_rule" "allow_all_traffic_ipv4" {
  security_group_id = aws_security_group.allow_ssh.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1"
}

resource "aws_instance" "main_server" {
  ami                    = var.ami
  instance_type          = "t2.micro"
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [aws_security_group.allow_ssh.id]
  key_name               = aws_key_pair.developer.key_name
  tags = {
    Name = "FB Main Server"
  }
}

resource "aws_key_pair" "developer" {
  key_name   = "thanhtd"
  public_key = var.public_key
}
