FROM docker.elastic.co/elasticsearch/elasticsearch:7.17.13

RUN yum install -y wget

RUN /usr/share/elasticsearch/bin/elasticsearch-plugin install --batch https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v7.17.7/elasticsearch-analysis-ik-7.17.7.zip

RUN /usr/share/elasticsearch/bin/elasticsearch-plugin install --batch https://github.com/medcl/elasticsearch-analysis-pinyin/releases/download/v7.17.7/elasticsearch-analysis-pinyin-7.17.7.zip
