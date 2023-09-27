# 创建用户索引
PUT user
{
  "settings": {
    "analysis": {
      "analyzer": {
        "text_analyzer": {
          "tokenizer": "ik_max_word",
          "filter": "py"
        },
        "completion_analyzer": {
          "tokenizer": "keyword",
          "filter": "py"
        }
      },
      "filter": {
        "py": {
          "type": "pinyin",
          "keep_full_pinyin": false,
          "keep_joined_full_pinyin": true,
          "keep_original": true,
          "limit_first_letter_length" : 16,
          "remove_duplicated_term" : true,
          "none_chinese_pinyin_tokenize": false
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "username": {
        "type": "text",
        "analyzer": "text_analyzer",
        "search_analyzer": "ik_max_word"
      },
      "follower_count": {
        "type": "integer"
      },
      "follow_count": {
        "type": "integer"
      },
      "suggestion": {
        "type": "completion",
        "analyzer": "completion_analyzer"
      },
      "create_time": {
        "type": "date",
        "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
      },
      "update_time": {
        "type": "date",
        "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
      },
      "version": {
        "type": "long"
      }
    }
  }
}

# 创建视频信息文档
PUT video
{
  "settings": {
    "analysis": {
      "analyzer": {
        "text_analyzer": {
          "tokenizer": "ik_max_word",
          "filter": "py"
        },
        "completion_analyzer": {
          "tokenizer": "keyword",
          "filter": "py"
        }
      },
      "filter": {
        "py": {
          "type": "pinyin",
          "keep_full_pinyin": false,
          "keep_joined_full_pinyin": true,
          "keep_original": true,
          "limit_first_letter_length" : 16,
          "remove_duplicated_term" : true,
          "none_chinese_pinyin_tokenize": false
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "title": {
        "type": "text",
        "analyzer": "text_analyzer",
        "search_analyzer": "ik_max_word"
      },
      "section_id": {
        "type": "keyword"
      },
      "tag_ids": {
        "type": "keyword"
      },
      "owner_id": {
        "type": "keyword"
      },
      "owner_name": {
        "type": "keyword"
      },
      "play_url": {
        "type": "keyword"
      },
      "cover_url": {
        "type": "keyword"
      },
      "comment_count": {
        "type": "long"
      },
      "like_count": {
        "type": "long"
      },
      "create_time": {
        "type": "date",
        "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
      },
      "update_time": {
        "type": "date",
        "format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
      },
      "suggestion": {
        "type": "completion",
        "analyzer": "completion_analyzer"
      },
      "version": {
        "type": "long"
      }
    }
  }
}