
resource "aws_security_group" "opensearch_sg" {
  vpc_id      = var.vpc_id
  name        = "opensearch_sg"
  description = "Allow Opensearch traffic"

  tags = {
    Name = "opensearch-sg"
  }
}

resource "aws_vpc_security_group_egress_rule" "opeansearch_egress" {
  security_group_id = aws_security_group.opensearch_sg.id
  from_port         = 0
  to_port           = 0
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"

}

resource "aws_vpc_security_group_ingress_rule" "opensearch_ingress" {
  security_group_id            = aws_security_group.opensearch_sg.id
  referenced_security_group_id = var.eks_node_sg_id
  from_port                    = 9200
  to_port                      = 9200
  ip_protocol                  = "tcp"
  cidr_ipv4                    = var.vpc_cidr
}

resource "aws_opensearch_domain" "aerona_opensearch" {
  domain_name    = "aerona-opensearch"
  engine_version = "Elasticsearch_7.10"
  cluster_config {
    instance_type  = "r4.large.search"
    instance_count = 2
  }
  vpc_options {
    subnet_ids         = [var.subnet_id]
    security_group_ids = [aws_security_group.opensearch_sg.id]
  }

  tags = {
    Domain = "DB Domain"
  }
}
