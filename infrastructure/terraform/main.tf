provider "aws" {
  region = "us-east-1"
}

# Create a VPC
resource "aws_vpc" "chaoslabs_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "ChaosLabsVPC"
  }
}

# Create a public subnet
resource "aws_subnet" "chaoslabs_subnet" {
  vpc_id                  = aws_vpc.chaoslabs_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true
  tags = {
    Name = "ChaosLabsSubnet"
  }
}

# Create an Internet Gateway
resource "aws_internet_gateway" "chaoslabs_igw" {
  vpc_id = aws_vpc.chaoslabs_vpc.id
  tags = {
    Name = "ChaosLabsIGW"
  }
}

# Create a Route Table and route to the Internet Gateway
resource "aws_route_table" "chaoslabs_rt" {
  vpc_id = aws_vpc.chaoslabs_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.chaoslabs_igw.id
  }

  tags = {
    Name = "ChaosLabsRouteTable"
  }
}

# Associate the Route Table with the Subnet
resource "aws_route_table_association" "chaoslabs_rt_assoc" {
  subnet_id      = aws_subnet.chaoslabs_subnet.id
  route_table_id = aws_route_table.chaoslabs_rt.id
}

# Launch an example EC2 instance (for demo)
resource "aws_instance" "chaoslabs_instance" {
  ami           = "ami-0c55b159cbfafe1f0"  # Example: Amazon Linux 2 AMI (us-east-1)
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.chaoslabs_subnet.id
  associate_public_ip_address = true
  tags = {
    Name = "ChaosLabsInstance"
  }
}
