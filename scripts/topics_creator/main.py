import logging
import sys
from typing import List

import ssl
import yaml

from configargparse import ArgumentParser, ArgumentDefaultsHelpFormatter, Namespace
from kafka.admin import KafkaAdminClient, NewTopic

SASL_MECHANISM = "PLAIN"
SASL_PLAINTEXT_SECURITY_PROTOCOL = "SASL_PLAINTEXT"
SASL_SSL_SECURITY_PROTOCOL = "SASL_SSL"

parser = ArgumentParser(
    auto_env_var_prefix="KAFKA_",
    formatter_class=ArgumentDefaultsHelpFormatter,
)
parser.add_argument("--brokers", type=str, required=True, default="127.0.0.1:9092",
                    help="Bootstrap servers separated by comma")
parser.add_argument("--user", type=str, required=True, help="SASL plain username")
parser.add_argument("--password", type=str, required=True, help="SASL plain password")
parser.add_argument("--ssl-cert", type=str, help="SSL certificate path")
parser.add_argument("--timeout_ms", type=int, default=5000, help="Topics create timeout in milliseconds")
parser.add_argument(
    "--log-level",
    type=str,
    choices=(
        logging.getLevelName(logging.DEBUG).lower(),
        logging.getLevelName(logging.INFO).lower(),
        logging.getLevelName(logging.WARNING).lower(),
        logging.getLevelName(logging.ERROR).lower(),
        logging.getLevelName(logging.CRITICAL).lower(),
    ),
    default=logging.getLevelName(logging.DEBUG).lower(),
    help="Log level"
)
parser.add_argument("topics_config_path", type=str, help="Path to the config for creating topics")


class TopicElement:
    __slots__ = ("topic", "partitions", "replicas")

    def __init__(self, topic: str, partitions: int = 1, replicas: int = 1) -> None:
        self.topic = topic
        self.partitions = partitions
        self.replicas = replicas

    def __repr__(self) -> str:
        return "{0}(topic='{1}', partitions={2}, replicas={3})".format(
            self.__class__.__name__,
            self.topic,
            self.partitions,
            self.replicas
        )


class TopicListElement:
    __slots__ = ("elements",)

    def __init__(self, arr: List[dict]) -> None:
        self.elements = []

        for obj in arr:
            self.elements.append(TopicElement(**obj))

    def __iter__(self):
        return iter(self.elements)


def main(args: Namespace) -> None:
    logger = logging.getLogger(__name__)
    logger.setLevel(args.log_level.upper())
    logger.addHandler(logging.StreamHandler(stream=sys.stdout))

    logger.info("Start to create kafka topics")
    logger.debug("args: %s", args)

    brokers = args.brokers.split(",")
    for i in range(len(brokers)):
        brokers[i] = brokers[i].strip()

    client_cfg = {
        "bootstrap_servers": brokers,
        "sasl_mechanism": SASL_MECHANISM,
        "security_protocol": SASL_PLAINTEXT_SECURITY_PROTOCOL,
        "sasl_plain_username": args.user,
        "sasl_plain_password": args.password,
    }

    if args.ssl_cert:
        logger.info("SSL mode is enabled")

        context = ssl.SSLContext(ssl.PROTOCOL_SSLv23)
        context.verify_mode = ssl.CERT_REQUIRED
        context.load_verify_locations(args.ssl_cert)

        client_cfg["ssl_context"] = context
        client_cfg["security_protocol"] = SASL_SSL_SECURITY_PROTOCOL

    logger.debug("kafka client config: %s", client_cfg)

    client = KafkaAdminClient(**client_cfg)
    exist_topic_names = set(client.list_topics())

    logger.debug("existing topics: %s", exist_topic_names)

    with open(args.topics_config_path, "r") as reader:
        topics_cfg = yaml.safe_load(reader)
        logger.debug("topics config: %s", topics_cfg)

    topics = []
    topic_elements = TopicListElement(topics_cfg)

    for element in topic_elements:
        if element.topic in exist_topic_names:
            continue

        # don't let set replicas more than number of brokers
        if element.replicas > len(brokers):
            element.replicas = len(brokers)

        topic = NewTopic(element.topic, element.partitions, element.replicas)
        topics.append(topic)

    if topics:
        logger.debug("topics to create: %s", [t.name for t in topics])
        client.create_topics(topics, args.timeout_ms, validate_only=False)
    else:
        logger.info("No need to create topics due everyone has already existed")

    client.close()
    logger.info("Topics creation completed successfully")


if __name__ == "__main__":
    main(parser.parse_args())
