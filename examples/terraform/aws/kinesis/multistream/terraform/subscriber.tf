module "lambda_subscriber" {
  source    = "../../../../../../build/terraform/aws/lambda"
  kms       = module.kms
  appconfig = module.appconfig

  config = {
    name        = "subscriber"
    description = "Substation node subscribes to Kinesis"
    image_uri   = "${module.ecr.url}:v1.2.0"
    image_arm   = true

    env = {
      "SUBSTATION_CONFIG" : "http://localhost:2772/applications/substation/environments/example/configurations/subscriber"
      "SUBSTATION_LAMBDA_HANDLER" : "AWS_KINESIS_DATA_STREAM"
      "SUBSTATION_DEBUG" : true
    }
  }

  depends_on = [
    module.appconfig.name,
    module.ecr.url,
  ]
}

resource "aws_lambda_event_source_mapping" "lambda_subscriber" {
  event_source_arn                   = module.kds_dst.arn
  function_name                      = module.lambda_subscriber.arn
  maximum_batching_window_in_seconds = 10
  batch_size                         = 100
  parallelization_factor             = 1
  starting_position                  = "LATEST"
}
