module "vpc" {
  source = "./modules/vpc"

  region = var.region
}

module "ec2" {
  source     = "./modules/ec2"
  vpc_id     = module.vpc.vpc_id
  subnet_id  = module.vpc.public_subnet_id
  public_key = var.public_key
}

module "elasticache" {
  source    = "./modules/elasticache"
  vpc_id    = module.vpc.vpc_id
  vpc_cidr  = module.vpc.vpc_cidr
  subnet_id = module.vpc.private_subnet_id
}

module "rds" {
  source      = "./modules/rds"
  vpc_cidr    = module.vpc.vpc_cidr
  vpc_id      = module.vpc.vpc_id
  subnet_id   = module.vpc.private_subnet_id
  db_password = var.postgres_password
}

module "opensearch" {
  source         = "./modules/opensearch"
  vpc_id         = module.vpc.private_subnet_id
  vpc_cidr       = module.vpc.vpc_cidr
  subnet_id      = module.vpc.private_subnet_id
  eks_node_sg_id = "" // TODO
}

module "documentdb" {
  source      = "./modules/docdb"
  vpc_id      = module.vpc.private_subnet_id
  vpc_cidr    = module.vpc.vpc_cidr
  subnet_id   = module.vpc.private_subnet_id
  db_password = var.docdb_password
}
