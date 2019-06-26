//resource "awsasgips_provider" "test-awsasgips" {
////
////  address = "1.2.3.6"
////
////}

//provider "aws" {
//  version = "~> 2.0"
//  region = "us-west-2"
//}

provider "awsasgips" {
  region = "us-west-2"

}

data "awsasgips" "test" {
  asgname = "asg-green-dev-media20190314181204859800000009"
  //region = "us-west-2"
}

//resource "null_resource" "testing-data-resource" {
//  //  triggers {
//  //    uuid = "${azurerm_virtual_machine.instance.id}"
//  //  }
//  provisioner "local-exec" {
////    when = "destroy"
//
//    command = <<EOT
//    echo "${data.awsasgips.test-awsasgips.output}"
//EOT
//  }
//}


output "instance_id" {
//  value = "${data.awsasgips.test.instance_id.0}"
  value = "${data.awsasgips.test.instance_id}"


}

output "private_ip" {
//  value = "${data.awsasgips.test.private_ip.0}"
  value = "${data.awsasgips.test.private_ip}"

}

output "private_ip_0" {
    value = "${data.awsasgips.test.private_ip.0}"

}
//output "public_ip" {
//  value = "${data.awsasgips.test.public_ip.0}"
//}


//Below is not working, map of list is not supported as output in data resources
//output "get_ip_from_output" {
//  value = "${data.awsasgips.test.output}"
//
//}