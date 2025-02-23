provider "aws" {
  region = "us-east-1"
}

# Create a VPC
resource "aws_vpc" "chaoslab_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "ChaosLabVPC"
  }
}

# Create a public subnet
resource "aws_subnet" "chaoslab_subnet" {
  vpc_id                  = aws_vpc.chaoslab_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true
  tags = {
    Name = "ChaosLabSubnet"
  }
}

# Create an Internet Gateway
resource "aws_internet_gateway" "chaoslab_igw" {
  vpc_id = aws_vpc.chaoslab_vpc.id
  tags = {
    Name = "ChaosLabIGW"
  }
}

# Create a Route Table and route to the Internet Gateway
resource "aws_route_table" "chaoslab_rt" {
  vpc_id = aws_vpc.chaoslab_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.chaoslab_igw.id
  }

  tags = {
    Name = "ChaosLabRouteTable"
  }
}

# Associate the Route Table with the Subnet
resource "aws_route_table_association" "chaoslab_rt_assoc" {
  subnet_id      = aws_subnet.chaoslab_subnet.id
  route_table_id = aws_route_table.chaoslab_rt.id
}

# Launch an example EC2 instance (for demonstration)
resource "aws_instance" "chaoslab_instance" {
  ami           = "ami-0c55b159cbfafe1f0"  # Example: Amazon Linux 2 AMI (us-east-1)
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.chaoslab_subnet.id
  associate_public_ip_address = true
  tags = {
    Name = "ChaosLabInstance"
  }
}
