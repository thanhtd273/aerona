
resource "aws_security_group" "redis_sg" {
  name        = "redis-sg"
  description = "Allow Redis traffic"
  vpc_id      = var.vpc_id

  tags = {
    Name = "redis-sg"
  }
}

resource "aws_vpc_security_group_egress_rule" "allow_redis_egress" {
  security_group_id = aws_security_group.redis_sg.id
  from_port         = 0
  to_port           = 0
  ip_protocol       = "-1"
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_vpc_security_group_ingress_rule" "allow_redis_ingress" {
  security_group_id = aws_security_group.redis_sg.id
  from_port         = 6379
  to_port           = 6379
  ip_protocol       = "tcp"
  cidr_ipv4         = var.vpc_cidr
}

resource "aws_elasticache_subnet_group" "redis" {
  name       = "aerona-redis-subnet"
  subnet_ids = [var.subnet_id]
}

resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "aerona-redis"
  engine               = "redis"
  engine_version       = var.engine_version
  node_type            = var.node_type
  num_cache_nodes      = 1
  parameter_group_name = var.parameter_group_name
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.redis.name
  security_group_ids   = [aws_security_group.redis_sg.id]

  tags = {
    Name = "aerona-redis"
  }
}

