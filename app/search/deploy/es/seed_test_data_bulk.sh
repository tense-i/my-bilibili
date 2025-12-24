#!/bin/bash
# seed_test_data_bulk.sh - 使用 _bulk API 批量写入测试数据
# 用法: ./seed_test_data_bulk.sh [ES_HOST]

set -e

# 取消所有代理设置
unset http_proxy
unset https_proxy
unset HTTP_PROXY
unset HTTPS_PROXY
unset all_proxy
unset ALL_PROXY

ES_HOST="${1:-http://localhost:9200}"

echo "ES 测试数据批量写入脚本"
echo "ES Host: $ES_HOST"
echo ""

# 检查 ES 连接
if ! curl -s "$ES_HOST" > /dev/null; then
    echo "错误: 无法连接到 ES: $ES_HOST"
    exit 1
fi

echo "1. 写入弹幕搜索数据 (dm_search_000)..."
cat > /tmp/dm_search_bulk.json << 'BULK_END'
{"index":{"_index":"dm_search_000","_id":"10001"}}
{"id":10001,"oid":1000,"oidstr":"1000","mid":1001,"content":"前方高能预警","mode":1,"pool":0,"progress":12000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16777215,"ctime":"2024-12-15 10:30:00","mtime":"2024-12-15 10:30:00"}
{"index":{"_index":"dm_search_000","_id":"10002"}}
{"id":10002,"oid":1000,"oidstr":"1000","mid":1002,"content":"太好看了","mode":1,"pool":0,"progress":15000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711680,"ctime":"2024-12-15 10:31:00","mtime":"2024-12-15 10:31:00"}
{"index":{"_index":"dm_search_000","_id":"10003"}}
{"id":10003,"oid":1000,"oidstr":"1000","mid":1003,"content":"泪目了","mode":1,"pool":0,"progress":18000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65280,"ctime":"2024-12-15 10:32:00","mtime":"2024-12-15 10:32:00"}
{"index":{"_index":"dm_search_000","_id":"10004"}}
{"id":10004,"oid":1000,"oidstr":"1000","mid":1004,"content":"哈哈哈哈哈","mode":1,"pool":0,"progress":20000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":255,"ctime":"2024-12-15 10:33:00","mtime":"2024-12-15 10:33:00"}
{"index":{"_index":"dm_search_000","_id":"10005"}}
{"id":10005,"oid":1000,"oidstr":"1000","mid":1005,"content":"这也太强了吧","mode":1,"pool":0,"progress":25000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16776960,"ctime":"2024-12-15 10:34:00","mtime":"2024-12-15 10:34:00"}
BULK_END
cat >> /tmp/dm_search_bulk.json << 'BULK_END'
{"index":{"_index":"dm_search_000","_id":"10006"}}
{"id":10006,"oid":1000,"oidstr":"1000","mid":1006,"content":"awsl","mode":1,"pool":0,"progress":28000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711935,"ctime":"2024-12-15 10:35:00","mtime":"2024-12-15 10:35:00"}
{"index":{"_index":"dm_search_000","_id":"10007"}}
{"id":10007,"oid":1000,"oidstr":"1000","mid":1007,"content":"名场面","mode":1,"pool":0,"progress":30000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65535,"ctime":"2024-12-15 10:36:00","mtime":"2024-12-15 10:36:00"}
{"index":{"_index":"dm_search_000","_id":"10008"}}
{"id":10008,"oid":1000,"oidstr":"1000","mid":1008,"content":"经典永流传","mode":1,"pool":0,"progress":32000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16777215,"ctime":"2024-12-15 10:37:00","mtime":"2024-12-15 10:37:00"}
{"index":{"_index":"dm_search_000","_id":"10009"}}
{"id":10009,"oid":1000,"oidstr":"1000","mid":1009,"content":"爷青回","mode":1,"pool":0,"progress":35000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711680,"ctime":"2024-12-15 10:38:00","mtime":"2024-12-15 10:38:00"}
{"index":{"_index":"dm_search_000","_id":"10010"}}
{"id":10010,"oid":1000,"oidstr":"1000","mid":1010,"content":"下次一定","mode":1,"pool":0,"progress":38000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65280,"ctime":"2024-12-15 10:39:00","mtime":"2024-12-15 10:39:00"}
{"index":{"_index":"dm_search_000","_id":"10011"}}
{"id":10011,"oid":1000,"oidstr":"1000","mid":1011,"content":"我来组成头部","mode":4,"pool":0,"progress":40000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":255,"ctime":"2024-12-15 10:40:00","mtime":"2024-12-15 10:40:00"}
{"index":{"_index":"dm_search_000","_id":"10012"}}
{"id":10012,"oid":1000,"oidstr":"1000","mid":1012,"content":"开幕雷击","mode":5,"pool":0,"progress":1000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16776960,"ctime":"2024-12-15 10:41:00","mtime":"2024-12-15 10:41:00"}
{"index":{"_index":"dm_search_000","_id":"10013"}}
{"id":10013,"oid":1000,"oidstr":"1000","mid":1013,"content":"标准结局","mode":1,"pool":0,"progress":290000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711935,"ctime":"2024-12-15 10:42:00","mtime":"2024-12-15 10:42:00"}
{"index":{"_index":"dm_search_000","_id":"10014"}}
{"id":10014,"oid":1000,"oidstr":"1000","mid":1014,"content":"这波啊这波是","mode":1,"pool":0,"progress":45000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65535,"ctime":"2024-12-15 10:43:00","mtime":"2024-12-15 10:43:00"}
{"index":{"_index":"dm_search_000","_id":"10015"}}
{"id":10015,"oid":1000,"oidstr":"1000","mid":1015,"content":"有内味了","mode":1,"pool":0,"progress":48000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16777215,"ctime":"2024-12-15 10:44:00","mtime":"2024-12-15 10:44:00"}
{"index":{"_index":"dm_search_000","_id":"10016"}}
{"id":10016,"oid":1000,"oidstr":"1000","mid":1016,"content":"好家伙","mode":1,"pool":0,"progress":50000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711680,"ctime":"2024-12-15 10:45:00","mtime":"2024-12-15 10:45:00"}
{"index":{"_index":"dm_search_000","_id":"10017"}}
{"id":10017,"oid":1000,"oidstr":"1000","mid":1017,"content":"绝绝子","mode":1,"pool":0,"progress":52000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65280,"ctime":"2024-12-15 10:46:00","mtime":"2024-12-15 10:46:00"}
{"index":{"_index":"dm_search_000","_id":"10018"}}
{"id":10018,"oid":1000,"oidstr":"1000","mid":1018,"content":"yyds","mode":1,"pool":0,"progress":55000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":255,"ctime":"2024-12-15 10:47:00","mtime":"2024-12-15 10:47:00"}
{"index":{"_index":"dm_search_000","_id":"10019"}}
{"id":10019,"oid":1000,"oidstr":"1000","mid":1019,"content":"破防了","mode":1,"pool":0,"progress":58000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16776960,"ctime":"2024-12-15 10:48:00","mtime":"2024-12-15 10:48:00"}
{"index":{"_index":"dm_search_000","_id":"10020"}}
{"id":10020,"oid":1000,"oidstr":"1000","mid":1020,"content":"DNA动了","mode":1,"pool":0,"progress":60000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711935,"ctime":"2024-12-15 10:49:00","mtime":"2024-12-15 10:49:00"}
BULK_END
cat >> /tmp/dm_search_bulk.json << 'BULK_END'
{"index":{"_index":"dm_search_000","_id":"10021"}}
{"id":10021,"oid":1000,"oidstr":"1000","mid":1021,"content":"前方核能","mode":1,"pool":0,"progress":62000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65535,"ctime":"2024-12-15 10:50:00","mtime":"2024-12-15 10:50:00"}
{"index":{"_index":"dm_search_000","_id":"10022"}}
{"id":10022,"oid":1000,"oidstr":"1000","mid":1022,"content":"太感动了","mode":1,"pool":0,"progress":65000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16777215,"ctime":"2024-12-15 10:51:00","mtime":"2024-12-15 10:51:00"}
{"index":{"_index":"dm_search_000","_id":"10023"}}
{"id":10023,"oid":1000,"oidstr":"1000","mid":1023,"content":"笑死我了","mode":1,"pool":0,"progress":68000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711680,"ctime":"2024-12-15 10:52:00","mtime":"2024-12-15 10:52:00"}
{"index":{"_index":"dm_search_000","_id":"10024"}}
{"id":10024,"oid":1000,"oidstr":"1000","mid":1024,"content":"这是什么神仙","mode":1,"pool":0,"progress":70000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65280,"ctime":"2024-12-15 10:53:00","mtime":"2024-12-15 10:53:00"}
{"index":{"_index":"dm_search_000","_id":"10025"}}
{"id":10025,"oid":1000,"oidstr":"1000","mid":1025,"content":"太强了","mode":1,"pool":0,"progress":72000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":255,"ctime":"2024-12-15 10:54:00","mtime":"2024-12-15 10:54:00"}
{"index":{"_index":"dm_search_000","_id":"10026"}}
{"id":10026,"oid":1000,"oidstr":"1000","mid":1026,"content":"这就是大佬吗","mode":1,"pool":0,"progress":75000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16776960,"ctime":"2024-12-15 10:55:00","mtime":"2024-12-15 10:55:00"}
{"index":{"_index":"dm_search_000","_id":"10027"}}
{"id":10027,"oid":1000,"oidstr":"1000","mid":1027,"content":"学到了","mode":1,"pool":0,"progress":78000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711935,"ctime":"2024-12-15 10:56:00","mtime":"2024-12-15 10:56:00"}
{"index":{"_index":"dm_search_000","_id":"10028"}}
{"id":10028,"oid":1000,"oidstr":"1000","mid":1028,"content":"涨知识了","mode":1,"pool":0,"progress":80000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65535,"ctime":"2024-12-15 10:57:00","mtime":"2024-12-15 10:57:00"}
{"index":{"_index":"dm_search_000","_id":"10029"}}
{"id":10029,"oid":1000,"oidstr":"1000","mid":1029,"content":"收藏了","mode":1,"pool":0,"progress":82000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16777215,"ctime":"2024-12-15 10:58:00","mtime":"2024-12-15 10:58:00"}
{"index":{"_index":"dm_search_000","_id":"10030"}}
{"id":10030,"oid":1000,"oidstr":"1000","mid":1030,"content":"投币了","mode":1,"pool":0,"progress":85000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711680,"ctime":"2024-12-15 10:59:00","mtime":"2024-12-15 10:59:00"}
{"index":{"_index":"dm_search_000","_id":"10031"}}
{"id":10031,"oid":1000,"oidstr":"1000","mid":1031,"content":"三连了","mode":1,"pool":0,"progress":88000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65280,"ctime":"2024-12-15 11:00:00","mtime":"2024-12-15 11:00:00"}
{"index":{"_index":"dm_search_000","_id":"10032"}}
{"id":10032,"oid":1000,"oidstr":"1000","mid":1032,"content":"催更","mode":1,"pool":0,"progress":90000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":255,"ctime":"2024-12-15 11:01:00","mtime":"2024-12-15 11:01:00"}
{"index":{"_index":"dm_search_000","_id":"10033"}}
{"id":10033,"oid":1000,"oidstr":"1000","mid":1033,"content":"期待下期","mode":1,"pool":0,"progress":92000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16776960,"ctime":"2024-12-15 11:02:00","mtime":"2024-12-15 11:02:00"}
{"index":{"_index":"dm_search_000","_id":"10034"}}
{"id":10034,"oid":1000,"oidstr":"1000","mid":1034,"content":"爱了爱了","mode":1,"pool":0,"progress":95000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711935,"ctime":"2024-12-15 11:03:00","mtime":"2024-12-15 11:03:00"}
{"index":{"_index":"dm_search_000","_id":"10035"}}
{"id":10035,"oid":1000,"oidstr":"1000","mid":1035,"content":"这波操作666","mode":1,"pool":0,"progress":98000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65535,"ctime":"2024-12-15 11:04:00","mtime":"2024-12-15 11:04:00"}
{"index":{"_index":"dm_search_000","_id":"10036"}}
{"id":10036,"oid":1000,"oidstr":"1000","mid":1036,"content":"神仙打架","mode":1,"pool":0,"progress":100000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16777215,"ctime":"2024-12-15 11:05:00","mtime":"2024-12-15 11:05:00"}
{"index":{"_index":"dm_search_000","_id":"10037"}}
{"id":10037,"oid":1000,"oidstr":"1000","mid":1037,"content":"太秀了","mode":1,"pool":0,"progress":102000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16711680,"ctime":"2024-12-15 11:06:00","mtime":"2024-12-15 11:06:00"}
{"index":{"_index":"dm_search_000","_id":"10038"}}
{"id":10038,"oid":1000,"oidstr":"1000","mid":1038,"content":"这就是天才吗","mode":1,"pool":0,"progress":105000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":65280,"ctime":"2024-12-15 11:07:00","mtime":"2024-12-15 11:07:00"}
{"index":{"_index":"dm_search_000","_id":"10039"}}
{"id":10039,"oid":1000,"oidstr":"1000","mid":1039,"content":"我悟了","mode":1,"pool":0,"progress":108000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":255,"ctime":"2024-12-15 11:08:00","mtime":"2024-12-15 11:08:00"}
{"index":{"_index":"dm_search_000","_id":"10040"}}
{"id":10040,"oid":1000,"oidstr":"1000","mid":1040,"content":"完结撒花","mode":1,"pool":0,"progress":295000,"state":0,"type":1,"attr":0,"attr_format":0,"fontsize":25,"color":16776960,"ctime":"2024-12-15 11:09:00","mtime":"2024-12-15 11:09:00"}
BULK_END

curl -s -X POST "$ES_HOST/_bulk" -H "Content-Type: application/x-ndjson" --data-binary @/tmp/dm_search_bulk.json > /dev/null
echo "  弹幕搜索数据写入完成 (40条)"

echo "2. 写入弹幕日期数据 (dm_date_2024_12)..."
cat > /tmp/dm_date_bulk.json << 'BULK_END'
{"index":{"_index":"dm_date_2024_12","_id":"20001"}}
{"id":20001,"oid":1001,"month":"2024-12","date":"2024-12-01","total":1523,"ctime":"2024-12-01 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20002"}}
{"id":20002,"oid":1002,"month":"2024-12","date":"2024-12-02","total":2341,"ctime":"2024-12-02 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20003"}}
{"id":20003,"oid":1003,"month":"2024-12","date":"2024-12-03","total":1876,"ctime":"2024-12-03 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20004"}}
{"id":20004,"oid":1004,"month":"2024-12","date":"2024-12-04","total":3245,"ctime":"2024-12-04 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20005"}}
{"id":20005,"oid":1005,"month":"2024-12","date":"2024-12-05","total":2567,"ctime":"2024-12-05 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20006"}}
{"id":20006,"oid":1006,"month":"2024-12","date":"2024-12-06","total":1934,"ctime":"2024-12-06 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20007"}}
{"id":20007,"oid":1007,"month":"2024-12","date":"2024-12-07","total":4521,"ctime":"2024-12-07 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20008"}}
{"id":20008,"oid":1008,"month":"2024-12","date":"2024-12-08","total":3876,"ctime":"2024-12-08 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20009"}}
{"id":20009,"oid":1009,"month":"2024-12","date":"2024-12-09","total":2134,"ctime":"2024-12-09 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20010"}}
{"id":20010,"oid":1010,"month":"2024-12","date":"2024-12-10","total":5678,"ctime":"2024-12-10 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20011"}}
{"id":20011,"oid":1011,"month":"2024-12","date":"2024-12-11","total":3421,"ctime":"2024-12-11 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20012"}}
{"id":20012,"oid":1012,"month":"2024-12","date":"2024-12-12","total":2987,"ctime":"2024-12-12 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20013"}}
{"id":20013,"oid":1013,"month":"2024-12","date":"2024-12-13","total":4123,"ctime":"2024-12-13 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20014"}}
{"id":20014,"oid":1014,"month":"2024-12","date":"2024-12-14","total":5234,"ctime":"2024-12-14 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20015"}}
{"id":20015,"oid":1015,"month":"2024-12","date":"2024-12-15","total":6789,"ctime":"2024-12-15 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20016"}}
{"id":20016,"oid":1016,"month":"2024-12","date":"2024-12-16","total":4567,"ctime":"2024-12-16 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20017"}}
{"id":20017,"oid":1017,"month":"2024-12","date":"2024-12-17","total":3890,"ctime":"2024-12-17 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20018"}}
{"id":20018,"oid":1018,"month":"2024-12","date":"2024-12-18","total":2456,"ctime":"2024-12-18 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20019"}}
{"id":20019,"oid":1019,"month":"2024-12","date":"2024-12-19","total":5123,"ctime":"2024-12-19 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20020"}}
{"id":20020,"oid":1020,"month":"2024-12","date":"2024-12-20","total":7890,"ctime":"2024-12-20 00:00:00"}
BULK_END
cat >> /tmp/dm_date_bulk.json << 'BULK_END'
{"index":{"_index":"dm_date_2024_12","_id":"20021"}}
{"id":20021,"oid":1021,"month":"2024-12","date":"2024-12-21","total":8234,"ctime":"2024-12-21 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20022"}}
{"id":20022,"oid":1022,"month":"2024-12","date":"2024-12-22","total":6543,"ctime":"2024-12-22 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20023"}}
{"id":20023,"oid":1023,"month":"2024-12","date":"2024-12-23","total":4321,"ctime":"2024-12-23 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20024"}}
{"id":20024,"oid":1024,"month":"2024-12","date":"2024-12-24","total":9876,"ctime":"2024-12-24 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20025"}}
{"id":20025,"oid":1025,"month":"2024-12","date":"2024-12-25","total":12345,"ctime":"2024-12-25 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20026"}}
{"id":20026,"oid":1026,"month":"2024-12","date":"2024-12-26","total":5678,"ctime":"2024-12-26 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20027"}}
{"id":20027,"oid":1027,"month":"2024-12","date":"2024-12-27","total":3456,"ctime":"2024-12-27 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20028"}}
{"id":20028,"oid":1028,"month":"2024-12","date":"2024-12-28","total":7654,"ctime":"2024-12-28 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20029"}}
{"id":20029,"oid":1029,"month":"2024-12","date":"2024-12-29","total":4567,"ctime":"2024-12-29 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20030"}}
{"id":20030,"oid":1030,"month":"2024-12","date":"2024-12-30","total":8901,"ctime":"2024-12-30 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20031"}}
{"id":20031,"oid":1031,"month":"2024-12","date":"2024-12-31","total":15678,"ctime":"2024-12-31 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20032"}}
{"id":20032,"oid":1000,"month":"2024-12","date":"2024-12-01","total":2345,"ctime":"2024-12-01 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20033"}}
{"id":20033,"oid":1000,"month":"2024-12","date":"2024-12-02","total":3456,"ctime":"2024-12-02 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20034"}}
{"id":20034,"oid":1000,"month":"2024-12","date":"2024-12-03","total":4567,"ctime":"2024-12-03 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20035"}}
{"id":20035,"oid":1000,"month":"2024-12","date":"2024-12-04","total":5678,"ctime":"2024-12-04 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20036"}}
{"id":20036,"oid":1000,"month":"2024-12","date":"2024-12-05","total":6789,"ctime":"2024-12-05 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20037"}}
{"id":20037,"oid":1000,"month":"2024-12","date":"2024-12-06","total":7890,"ctime":"2024-12-06 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20038"}}
{"id":20038,"oid":1000,"month":"2024-12","date":"2024-12-07","total":8901,"ctime":"2024-12-07 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20039"}}
{"id":20039,"oid":1000,"month":"2024-12","date":"2024-12-08","total":9012,"ctime":"2024-12-08 00:00:00"}
{"index":{"_index":"dm_date_2024_12","_id":"20040"}}
{"id":20040,"oid":1000,"month":"2024-12","date":"2024-12-09","total":10123,"ctime":"2024-12-09 00:00:00"}
BULK_END

curl -s -X POST "$ES_HOST/_bulk" -H "Content-Type: application/x-ndjson" --data-binary @/tmp/dm_date_bulk.json > /dev/null
echo "  弹幕日期数据写入完成 (40条)"

echo "3. 写入 PGC 番剧数据 (pgc_media)..."
cat > /tmp/pgc_media_bulk.json << 'BULK_END'
{"index":{"_index":"pgc_media","_id":"28220001"}}
{"media_id":28220001,"season_id":39001,"title":"进击的巨人 最终季","season_type":1,"style_id":1,"status":0,"release_date":"2024-01-15","producer_id":101,"is_deleted":0,"area_id":"1","score":9.8,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-01-15 00:00:00","season_month":1,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":5000000,"play_count":100000000,"fav_count":5000000,"ctime":"2024-01-15 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220002"}}
{"media_id":28220002,"season_id":39002,"title":"鬼灭之刃 无限列车篇","season_type":1,"style_id":2,"status":0,"release_date":"2024-02-20","producer_id":102,"is_deleted":0,"area_id":"1","score":9.7,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-02-20 00:00:00","season_month":2,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":4500000,"play_count":95000000,"fav_count":4800000,"ctime":"2024-02-20 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220003"}}
{"media_id":28220003,"season_id":39003,"title":"咒术回战 第二季","season_type":1,"style_id":3,"status":0,"release_date":"2024-03-10","producer_id":103,"is_deleted":0,"area_id":"1","score":9.6,"is_finish":"0","season_version":2,"season_status":1,"pub_time":"2024-03-10 00:00:00","season_month":3,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":4000000,"play_count":85000000,"fav_count":4200000,"ctime":"2024-03-10 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220004"}}
{"media_id":28220004,"season_id":39004,"title":"间谍过家家 第二季","season_type":1,"style_id":4,"status":0,"release_date":"2024-04-05","producer_id":104,"is_deleted":0,"area_id":"1","score":9.5,"is_finish":"0","season_version":2,"season_status":1,"pub_time":"2024-04-05 00:00:00","season_month":4,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3800000,"play_count":80000000,"fav_count":4000000,"ctime":"2024-04-05 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220005"}}
{"media_id":28220005,"season_id":39005,"title":"电锯人","season_type":1,"style_id":5,"status":0,"release_date":"2024-05-15","producer_id":105,"is_deleted":0,"area_id":"1","score":9.4,"is_finish":"0","season_version":1,"season_status":1,"pub_time":"2024-05-15 00:00:00","season_month":5,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3500000,"play_count":75000000,"fav_count":3800000,"ctime":"2024-05-15 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220006"}}
{"media_id":28220006,"season_id":39006,"title":"我的英雄学院 第七季","season_type":1,"style_id":6,"status":0,"release_date":"2024-06-01","producer_id":106,"is_deleted":0,"area_id":"1","score":9.3,"is_finish":"0","season_version":7,"season_status":1,"pub_time":"2024-06-01 00:00:00","season_month":6,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3200000,"play_count":70000000,"fav_count":3500000,"ctime":"2024-06-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220007"}}
{"media_id":28220007,"season_id":39007,"title":"JOJO的奇妙冒险 石之海","season_type":1,"style_id":7,"status":0,"release_date":"2024-07-10","producer_id":107,"is_deleted":0,"area_id":"1","score":9.5,"is_finish":"1","season_version":6,"season_status":0,"pub_time":"2024-07-10 00:00:00","season_month":7,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2800000,"play_count":65000000,"fav_count":3200000,"ctime":"2024-07-10 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220008"}}
{"media_id":28220008,"season_id":39008,"title":"辉夜大小姐想让我告白 第三季","season_type":1,"style_id":8,"status":0,"release_date":"2024-08-20","producer_id":108,"is_deleted":0,"area_id":"1","score":9.6,"is_finish":"1","season_version":3,"season_status":0,"pub_time":"2024-08-20 00:00:00","season_month":8,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2500000,"play_count":60000000,"fav_count":3000000,"ctime":"2024-08-20 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220009"}}
{"media_id":28220009,"season_id":39009,"title":"孤独摇滚","season_type":1,"style_id":9,"status":0,"release_date":"2024-09-15","producer_id":109,"is_deleted":0,"area_id":"1","score":9.7,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-09-15 00:00:00","season_month":9,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2200000,"play_count":55000000,"fav_count":2800000,"ctime":"2024-09-15 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220010"}}
{"media_id":28220010,"season_id":39010,"title":"葬送的芙莉莲","season_type":1,"style_id":10,"status":0,"release_date":"2024-10-01","producer_id":110,"is_deleted":0,"area_id":"1","score":9.9,"is_finish":"0","season_version":1,"season_status":1,"pub_time":"2024-10-01 00:00:00","season_month":10,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":6000000,"play_count":120000000,"fav_count":6500000,"ctime":"2024-10-01 00:00:00","mtime":"2024-12-15 00:00:00"}
BULK_END
cat >> /tmp/pgc_media_bulk.json << 'BULK_END'
{"index":{"_index":"pgc_media","_id":"28220011"}}
{"media_id":28220011,"season_id":39011,"title":"药屋少女的呢喃","season_type":1,"style_id":11,"status":0,"release_date":"2024-10-15","producer_id":111,"is_deleted":0,"area_id":"1","score":9.5,"is_finish":"0","season_version":1,"season_status":1,"pub_time":"2024-10-15 00:00:00","season_month":10,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3500000,"play_count":75000000,"fav_count":3800000,"ctime":"2024-10-15 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220012"}}
{"media_id":28220012,"season_id":39012,"title":"迷宫饭","season_type":1,"style_id":12,"status":0,"release_date":"2024-11-01","producer_id":112,"is_deleted":0,"area_id":"1","score":9.4,"is_finish":"0","season_version":1,"season_status":1,"pub_time":"2024-11-01 00:00:00","season_month":11,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2800000,"play_count":60000000,"fav_count":3000000,"ctime":"2024-11-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220013"}}
{"media_id":28220013,"season_id":39013,"title":"怪兽8号","season_type":1,"style_id":13,"status":0,"release_date":"2024-11-15","producer_id":113,"is_deleted":0,"area_id":"1","score":9.2,"is_finish":"0","season_version":1,"season_status":1,"pub_time":"2024-11-15 00:00:00","season_month":11,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2500000,"play_count":55000000,"fav_count":2700000,"ctime":"2024-11-15 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220014"}}
{"media_id":28220014,"season_id":39014,"title":"排球少年 垃圾场的决战","season_type":2,"style_id":14,"status":0,"release_date":"2024-12-01","producer_id":114,"is_deleted":0,"area_id":"1","score":9.8,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-12-01 00:00:00","season_month":12,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":4000000,"play_count":85000000,"fav_count":4200000,"ctime":"2024-12-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220015"}}
{"media_id":28220015,"season_id":39015,"title":"蓝色监狱 第二季","season_type":1,"style_id":15,"status":0,"release_date":"2024-12-10","producer_id":115,"is_deleted":0,"area_id":"1","score":9.3,"is_finish":"0","season_version":2,"season_status":1,"pub_time":"2024-12-10 00:00:00","season_month":12,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3000000,"play_count":65000000,"fav_count":3200000,"ctime":"2024-12-10 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220016"}}
{"media_id":28220016,"season_id":39016,"title":"天国大魔境","season_type":1,"style_id":16,"status":0,"release_date":"2024-03-20","producer_id":116,"is_deleted":0,"area_id":"1","score":9.4,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-03-20 00:00:00","season_month":3,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2200000,"play_count":50000000,"fav_count":2500000,"ctime":"2024-03-20 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220017"}}
{"media_id":28220017,"season_id":39017,"title":"地狱乐","season_type":1,"style_id":17,"status":0,"release_date":"2024-04-15","producer_id":117,"is_deleted":0,"area_id":"1","score":9.1,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-04-15 00:00:00","season_month":4,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2000000,"play_count":45000000,"fav_count":2200000,"ctime":"2024-04-15 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220018"}}
{"media_id":28220018,"season_id":39018,"title":"我推的孩子","season_type":1,"style_id":18,"status":0,"release_date":"2024-05-01","producer_id":118,"is_deleted":0,"area_id":"1","score":9.6,"is_finish":"0","season_version":1,"season_status":1,"pub_time":"2024-05-01 00:00:00","season_month":5,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":4500000,"play_count":95000000,"fav_count":4800000,"ctime":"2024-05-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220019"}}
{"media_id":28220019,"season_id":39019,"title":"无职转生 第二季","season_type":1,"style_id":19,"status":0,"release_date":"2024-06-15","producer_id":119,"is_deleted":0,"area_id":"1","score":9.5,"is_finish":"0","season_version":2,"season_status":1,"pub_time":"2024-06-15 00:00:00","season_month":6,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3800000,"play_count":80000000,"fav_count":4000000,"ctime":"2024-06-15 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220020"}}
{"media_id":28220020,"season_id":39020,"title":"86-不存在的战区- 第二季","season_type":1,"style_id":20,"status":0,"release_date":"2024-07-01","producer_id":120,"is_deleted":0,"area_id":"1","score":9.4,"is_finish":"1","season_version":2,"season_status":0,"pub_time":"2024-07-01 00:00:00","season_month":7,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2500000,"play_count":55000000,"fav_count":2800000,"ctime":"2024-07-01 00:00:00","mtime":"2024-12-15 00:00:00"}
BULK_END
cat >> /tmp/pgc_media_bulk.json << 'BULK_END'
{"index":{"_index":"pgc_media","_id":"28220021"}}
{"media_id":28220021,"season_id":39021,"title":"刀剑神域 进击篇","season_type":1,"style_id":21,"status":0,"release_date":"2024-08-01","producer_id":121,"is_deleted":0,"area_id":"1","score":9.2,"is_finish":"0","season_version":1,"season_status":1,"pub_time":"2024-08-01 00:00:00","season_month":8,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3200000,"play_count":70000000,"fav_count":3500000,"ctime":"2024-08-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220022"}}
{"media_id":28220022,"season_id":39022,"title":"关于我转生变成史莱姆这档事 第三季","season_type":1,"style_id":22,"status":0,"release_date":"2024-09-01","producer_id":122,"is_deleted":0,"area_id":"1","score":9.3,"is_finish":"0","season_version":3,"season_status":1,"pub_time":"2024-09-01 00:00:00","season_month":9,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3500000,"play_count":75000000,"fav_count":3800000,"ctime":"2024-09-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220023"}}
{"media_id":28220023,"season_id":39023,"title":"盾之勇者成名录 第三季","season_type":1,"style_id":23,"status":0,"release_date":"2024-10-10","producer_id":123,"is_deleted":0,"area_id":"1","score":8.9,"is_finish":"0","season_version":3,"season_status":1,"pub_time":"2024-10-10 00:00:00","season_month":10,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2000000,"play_count":45000000,"fav_count":2200000,"ctime":"2024-10-10 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220024"}}
{"media_id":28220024,"season_id":39024,"title":"Re:从零开始的异世界生活 第三季","season_type":1,"style_id":24,"status":0,"release_date":"2024-11-10","producer_id":124,"is_deleted":0,"area_id":"1","score":9.5,"is_finish":"0","season_version":3,"season_status":1,"pub_time":"2024-11-10 00:00:00","season_month":11,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":4000000,"play_count":85000000,"fav_count":4200000,"ctime":"2024-11-10 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220025"}}
{"media_id":28220025,"season_id":39025,"title":"overlord 第五季","season_type":1,"style_id":25,"status":0,"release_date":"2024-12-05","producer_id":125,"is_deleted":0,"area_id":"1","score":9.4,"is_finish":"0","season_version":5,"season_status":1,"pub_time":"2024-12-05 00:00:00","season_month":12,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3000000,"play_count":65000000,"fav_count":3200000,"ctime":"2024-12-05 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220026"}}
{"media_id":28220026,"season_id":39026,"title":"你的名字","season_type":2,"style_id":26,"status":0,"release_date":"2024-01-01","producer_id":126,"is_deleted":0,"area_id":"1","score":9.9,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-01-01 00:00:00","season_month":1,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":8000000,"play_count":200000000,"fav_count":10000000,"ctime":"2024-01-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220027"}}
{"media_id":28220027,"season_id":39027,"title":"铃芽之旅","season_type":2,"style_id":27,"status":0,"release_date":"2024-02-01","producer_id":127,"is_deleted":0,"area_id":"1","score":9.7,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-02-01 00:00:00","season_month":2,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":6000000,"play_count":150000000,"fav_count":7500000,"ctime":"2024-02-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220028"}}
{"media_id":28220028,"season_id":39028,"title":"灌篮高手 THE FIRST SLAM DUNK","season_type":2,"style_id":28,"status":0,"release_date":"2024-03-01","producer_id":128,"is_deleted":0,"area_id":"1","score":9.8,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-03-01 00:00:00","season_month":3,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":7000000,"play_count":180000000,"fav_count":9000000,"ctime":"2024-03-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220029"}}
{"media_id":28220029,"season_id":39029,"title":"工作细胞 剧场版","season_type":2,"style_id":29,"status":0,"release_date":"2024-04-01","producer_id":129,"is_deleted":0,"area_id":"1","score":9.3,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-04-01 00:00:00","season_month":4,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3000000,"play_count":70000000,"fav_count":3500000,"ctime":"2024-04-01 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220030"}}
{"media_id":28220030,"season_id":39030,"title":"我想吃掉你的胰脏","season_type":2,"style_id":30,"status":0,"release_date":"2024-05-10","producer_id":130,"is_deleted":0,"area_id":"1","score":9.5,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-05-10 00:00:00","season_month":5,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":4000000,"play_count":90000000,"fav_count":4500000,"ctime":"2024-05-10 00:00:00","mtime":"2024-12-15 00:00:00"}
BULK_END
cat >> /tmp/pgc_media_bulk.json << 'BULK_END'
{"index":{"_index":"pgc_media","_id":"28220031"}}
{"media_id":28220031,"season_id":39031,"title":"声之形","season_type":2,"style_id":31,"status":0,"release_date":"2024-06-10","producer_id":131,"is_deleted":0,"area_id":"1","score":9.6,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-06-10 00:00:00","season_month":6,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":5000000,"play_count":110000000,"fav_count":5500000,"ctime":"2024-06-10 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220032"}}
{"media_id":28220032,"season_id":39032,"title":"紫罗兰永恒花园 剧场版","season_type":2,"style_id":32,"status":0,"release_date":"2024-07-15","producer_id":132,"is_deleted":0,"area_id":"1","score":9.8,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-07-15 00:00:00","season_month":7,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":6500000,"play_count":140000000,"fav_count":7000000,"ctime":"2024-07-15 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220033"}}
{"media_id":28220033,"season_id":39033,"title":"我们仍未知道那天所看见的花的名字 剧场版","season_type":2,"style_id":33,"status":0,"release_date":"2024-08-10","producer_id":133,"is_deleted":0,"area_id":"1","score":9.4,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-08-10 00:00:00","season_month":8,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":4500000,"play_count":100000000,"fav_count":5000000,"ctime":"2024-08-10 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220034"}}
{"media_id":28220034,"season_id":39034,"title":"天气之子","season_type":2,"style_id":34,"status":0,"release_date":"2024-09-10","producer_id":134,"is_deleted":0,"area_id":"1","score":9.5,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-09-10 00:00:00","season_month":9,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":5500000,"play_count":130000000,"fav_count":6500000,"ctime":"2024-09-10 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220035"}}
{"media_id":28220035,"season_id":39035,"title":"言叶之庭","season_type":2,"style_id":35,"status":0,"release_date":"2024-10-20","producer_id":135,"is_deleted":0,"area_id":"1","score":9.3,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-10-20 00:00:00","season_month":10,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3500000,"play_count":80000000,"fav_count":4000000,"ctime":"2024-10-20 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220036"}}
{"media_id":28220036,"season_id":39036,"title":"秒速五厘米","season_type":2,"style_id":36,"status":0,"release_date":"2024-11-20","producer_id":136,"is_deleted":0,"area_id":"1","score":9.2,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-11-20 00:00:00","season_month":11,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3000000,"play_count":70000000,"fav_count":3500000,"ctime":"2024-11-20 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220037"}}
{"media_id":28220037,"season_id":39037,"title":"夏日重现","season_type":1,"style_id":37,"status":0,"release_date":"2024-04-20","producer_id":137,"is_deleted":0,"area_id":"1","score":9.6,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-04-20 00:00:00","season_month":4,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":4200000,"play_count":92000000,"fav_count":4600000,"ctime":"2024-04-20 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220038"}}
{"media_id":28220038,"season_id":39038,"title":"赛博朋克 边缘行者","season_type":1,"style_id":38,"status":0,"release_date":"2024-05-20","producer_id":138,"is_deleted":0,"area_id":"1","score":9.7,"is_finish":"1","season_version":1,"season_status":0,"pub_time":"2024-05-20 00:00:00","season_month":5,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":5200000,"play_count":115000000,"fav_count":5700000,"ctime":"2024-05-20 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220039"}}
{"media_id":28220039,"season_id":39039,"title":"石纪元 第三季","season_type":1,"style_id":39,"status":0,"release_date":"2024-06-20","producer_id":139,"is_deleted":0,"area_id":"1","score":9.3,"is_finish":"0","season_version":3,"season_status":1,"pub_time":"2024-06-20 00:00:00","season_month":6,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":2800000,"play_count":62000000,"fav_count":3100000,"ctime":"2024-06-20 00:00:00","mtime":"2024-12-15 00:00:00"}
{"index":{"_index":"pgc_media","_id":"28220040"}}
{"media_id":28220040,"season_id":39040,"title":"国王排名 第二季","season_type":1,"style_id":40,"status":0,"release_date":"2024-07-20","producer_id":140,"is_deleted":0,"area_id":"1","score":9.5,"is_finish":"0","season_version":2,"season_status":1,"pub_time":"2024-07-20 00:00:00","season_month":7,"latest_time":"2024-12-15 00:00:00","copyright_info":"bilibili","dm_count":3300000,"play_count":72000000,"fav_count":3600000,"ctime":"2024-07-20 00:00:00","mtime":"2024-12-15 00:00:00"}
BULK_END

curl -s -X POST "$ES_HOST/_bulk" -H "Content-Type: application/x-ndjson" --data-binary @/tmp/pgc_media_bulk.json > /dev/null
echo "  PGC 番剧数据写入完成 (40条)"

echo "4. 写入评论记录数据 (replyrecord_00)..."
cat > /tmp/replyrecord_bulk.json << 'BULK_END'
{"index":{"_index":"replyrecord_00","_id":"30001_1001"}}
{"id":30001,"oid":1001,"mid":100,"type":1,"state":0,"content":"这个视频太棒了！","like":1523,"hate":12,"rcount":45,"floor":1,"ctime":"2024-12-15 10:30:00","mtime":"2024-12-15 10:30:00"}
{"index":{"_index":"replyrecord_00","_id":"30002_1002"}}
{"id":30002,"oid":1002,"mid":100,"type":1,"state":0,"content":"UP主辛苦了","like":2341,"hate":8,"rcount":67,"floor":2,"ctime":"2024-12-15 10:31:00","mtime":"2024-12-15 10:31:00"}
{"index":{"_index":"replyrecord_00","_id":"30003_1003"}}
{"id":30003,"oid":1003,"mid":100,"type":1,"state":0,"content":"三连支持","like":1876,"hate":5,"rcount":34,"floor":3,"ctime":"2024-12-15 10:32:00","mtime":"2024-12-15 10:32:00"}
{"index":{"_index":"replyrecord_00","_id":"30004_1004"}}
{"id":30004,"oid":1004,"mid":100,"type":1,"state":0,"content":"学到了","like":3245,"hate":15,"rcount":89,"floor":4,"ctime":"2024-12-15 10:33:00","mtime":"2024-12-15 10:33:00"}
{"index":{"_index":"replyrecord_00","_id":"30005_1005"}}
{"id":30005,"oid":1005,"mid":100,"type":1,"state":0,"content":"感谢分享","like":2567,"hate":10,"rcount":56,"floor":5,"ctime":"2024-12-15 10:34:00","mtime":"2024-12-15 10:34:00"}
{"index":{"_index":"replyrecord_00","_id":"30006_1006"}}
{"id":30006,"oid":1006,"mid":100,"type":1,"state":0,"content":"催更催更","like":1934,"hate":7,"rcount":42,"floor":6,"ctime":"2024-12-15 10:35:00","mtime":"2024-12-15 10:35:00"}
{"index":{"_index":"replyrecord_00","_id":"30007_1007"}}
{"id":30007,"oid":1007,"mid":100,"type":1,"state":0,"content":"这期质量很高","like":4521,"hate":20,"rcount":123,"floor":7,"ctime":"2024-12-15 10:36:00","mtime":"2024-12-15 10:36:00"}
{"index":{"_index":"replyrecord_00","_id":"30008_1008"}}
{"id":30008,"oid":1008,"mid":100,"type":1,"state":0,"content":"看完了，意犹未尽","like":3876,"hate":18,"rcount":98,"floor":8,"ctime":"2024-12-15 10:37:00","mtime":"2024-12-15 10:37:00"}
{"index":{"_index":"replyrecord_00","_id":"30009_1009"}}
{"id":30009,"oid":1009,"mid":100,"type":1,"state":0,"content":"建议收藏","like":2134,"hate":9,"rcount":51,"floor":9,"ctime":"2024-12-15 10:38:00","mtime":"2024-12-15 10:38:00"}
{"index":{"_index":"replyrecord_00","_id":"30010_1010"}}
{"id":30010,"oid":1010,"mid":100,"type":1,"state":0,"content":"已投币","like":5678,"hate":25,"rcount":156,"floor":10,"ctime":"2024-12-15 10:39:00","mtime":"2024-12-15 10:39:00"}
{"index":{"_index":"replyrecord_00","_id":"30011_1011"}}
{"id":30011,"oid":1011,"mid":100,"type":1,"state":0,"content":"前排占座","like":3421,"hate":14,"rcount":78,"floor":11,"ctime":"2024-12-15 10:40:00","mtime":"2024-12-15 10:40:00"}
{"index":{"_index":"replyrecord_00","_id":"30012_1012"}}
{"id":30012,"oid":1012,"mid":100,"type":1,"state":0,"content":"来晚了","like":2987,"hate":11,"rcount":65,"floor":12,"ctime":"2024-12-15 10:41:00","mtime":"2024-12-15 10:41:00"}
{"index":{"_index":"replyrecord_00","_id":"30013_1013"}}
{"id":30013,"oid":1013,"mid":100,"type":1,"state":0,"content":"每日打卡","like":4123,"hate":19,"rcount":112,"floor":13,"ctime":"2024-12-15 10:42:00","mtime":"2024-12-15 10:42:00"}
{"index":{"_index":"replyrecord_00","_id":"30014_1014"}}
{"id":30014,"oid":1014,"mid":100,"type":1,"state":0,"content":"这就是我想看的","like":5234,"hate":23,"rcount":145,"floor":14,"ctime":"2024-12-15 10:43:00","mtime":"2024-12-15 10:43:00"}
{"index":{"_index":"replyrecord_00","_id":"30015_1015"}}
{"id":30015,"oid":1015,"mid":100,"type":1,"state":0,"content":"太强了","like":6789,"hate":30,"rcount":189,"floor":15,"ctime":"2024-12-15 10:44:00","mtime":"2024-12-15 10:44:00"}
{"index":{"_index":"replyrecord_00","_id":"30016_1016"}}
{"id":30016,"oid":1016,"mid":100,"type":1,"state":0,"content":"涨知识了","like":4567,"hate":21,"rcount":134,"floor":16,"ctime":"2024-12-15 10:45:00","mtime":"2024-12-15 10:45:00"}
{"index":{"_index":"replyrecord_00","_id":"30017_1017"}}
{"id":30017,"oid":1017,"mid":100,"type":1,"state":0,"content":"期待下期","like":3890,"hate":17,"rcount":101,"floor":17,"ctime":"2024-12-15 10:46:00","mtime":"2024-12-15 10:46:00"}
{"index":{"_index":"replyrecord_00","_id":"30018_1018"}}
{"id":30018,"oid":1018,"mid":100,"type":1,"state":0,"content":"收藏从未停止","like":2456,"hate":10,"rcount":58,"floor":18,"ctime":"2024-12-15 10:47:00","mtime":"2024-12-15 10:47:00"}
{"index":{"_index":"replyrecord_00","_id":"30019_1019"}}
{"id":30019,"oid":1019,"mid":100,"type":1,"state":0,"content":"这波操作666","like":5123,"hate":22,"rcount":142,"floor":19,"ctime":"2024-12-15 10:48:00","mtime":"2024-12-15 10:48:00"}
{"index":{"_index":"replyrecord_00","_id":"30020_1020"}}
{"id":30020,"oid":1020,"mid":100,"type":1,"state":0,"content":"爱了爱了","like":7890,"hate":35,"rcount":212,"floor":20,"ctime":"2024-12-15 10:49:00","mtime":"2024-12-15 10:49:00"}
BULK_END
cat >> /tmp/replyrecord_bulk.json << 'BULK_END'
{"index":{"_index":"replyrecord_00","_id":"30021_1021"}}
{"id":30021,"oid":1021,"mid":100,"type":1,"state":0,"content":"神仙UP主","like":8234,"hate":38,"rcount":234,"floor":21,"ctime":"2024-12-15 10:50:00","mtime":"2024-12-15 10:50:00"}
{"index":{"_index":"replyrecord_00","_id":"30022_1022"}}
{"id":30022,"oid":1022,"mid":100,"type":1,"state":0,"content":"这就是专业","like":6543,"hate":28,"rcount":178,"floor":22,"ctime":"2024-12-15 10:51:00","mtime":"2024-12-15 10:51:00"}
{"index":{"_index":"replyrecord_00","_id":"30023_1023"}}
{"id":30023,"oid":1023,"mid":100,"type":1,"state":0,"content":"太秀了","like":4321,"hate":19,"rcount":118,"floor":23,"ctime":"2024-12-15 10:52:00","mtime":"2024-12-15 10:52:00"}
{"index":{"_index":"replyrecord_00","_id":"30024_1024"}}
{"id":30024,"oid":1024,"mid":100,"type":1,"state":0,"content":"这是什么神仙视频","like":9876,"hate":45,"rcount":267,"floor":24,"ctime":"2024-12-15 10:53:00","mtime":"2024-12-15 10:53:00"}
{"index":{"_index":"replyrecord_00","_id":"30025_1025"}}
{"id":30025,"oid":1025,"mid":100,"type":1,"state":0,"content":"我悟了","like":12345,"hate":55,"rcount":345,"floor":25,"ctime":"2024-12-15 10:54:00","mtime":"2024-12-15 10:54:00"}
{"index":{"_index":"replyrecord_00","_id":"30026_1026"}}
{"id":30026,"oid":1026,"mid":100,"type":2,"state":0,"content":"话题讨论：这个观点很有意思","like":5678,"hate":24,"rcount":156,"floor":26,"ctime":"2024-12-15 10:55:00","mtime":"2024-12-15 10:55:00"}
{"index":{"_index":"replyrecord_00","_id":"30027_1027"}}
{"id":30027,"oid":1027,"mid":100,"type":2,"state":0,"content":"同意楼上的看法","like":3456,"hate":15,"rcount":89,"floor":27,"ctime":"2024-12-15 10:56:00","mtime":"2024-12-15 10:56:00"}
{"index":{"_index":"replyrecord_00","_id":"30028_1028"}}
{"id":30028,"oid":1028,"mid":100,"type":2,"state":0,"content":"我有不同意见","like":7654,"hate":32,"rcount":201,"floor":28,"ctime":"2024-12-15 10:57:00","mtime":"2024-12-15 10:57:00"}
{"index":{"_index":"replyrecord_00","_id":"30029_1029"}}
{"id":30029,"oid":1029,"mid":100,"type":2,"state":0,"content":"这个话题很有深度","like":4567,"hate":20,"rcount":123,"floor":29,"ctime":"2024-12-15 10:58:00","mtime":"2024-12-15 10:58:00"}
{"index":{"_index":"replyrecord_00","_id":"30030_1030"}}
{"id":30030,"oid":1030,"mid":100,"type":2,"state":0,"content":"期待更多讨论","like":8901,"hate":40,"rcount":245,"floor":30,"ctime":"2024-12-15 10:59:00","mtime":"2024-12-15 10:59:00"}
{"index":{"_index":"replyrecord_00","_id":"30031_1031"}}
{"id":30031,"oid":1031,"mid":100,"type":3,"state":0,"content":"活动参与：支持这个活动","like":15678,"hate":70,"rcount":423,"floor":31,"ctime":"2024-12-15 11:00:00","mtime":"2024-12-15 11:00:00"}
{"index":{"_index":"replyrecord_00","_id":"30032_1032"}}
{"id":30032,"oid":1032,"mid":100,"type":3,"state":0,"content":"活动太棒了","like":2345,"hate":10,"rcount":62,"floor":32,"ctime":"2024-12-15 11:01:00","mtime":"2024-12-15 11:01:00"}
{"index":{"_index":"replyrecord_00","_id":"30033_1033"}}
{"id":30033,"oid":1033,"mid":100,"type":3,"state":0,"content":"已参与活动","like":3456,"hate":14,"rcount":91,"floor":33,"ctime":"2024-12-15 11:02:00","mtime":"2024-12-15 11:02:00"}
{"index":{"_index":"replyrecord_00","_id":"30034_1034"}}
{"id":30034,"oid":1034,"mid":100,"type":3,"state":0,"content":"活动奖品很丰富","like":4567,"hate":19,"rcount":124,"floor":34,"ctime":"2024-12-15 11:03:00","mtime":"2024-12-15 11:03:00"}
{"index":{"_index":"replyrecord_00","_id":"30035_1035"}}
{"id":30035,"oid":1035,"mid":100,"type":3,"state":0,"content":"希望能中奖","like":5678,"hate":24,"rcount":156,"floor":35,"ctime":"2024-12-15 11:04:00","mtime":"2024-12-15 11:04:00"}
{"index":{"_index":"replyrecord_00","_id":"30036_1036"}}
{"id":30036,"oid":1036,"mid":100,"type":1,"state":0,"content":"这个系列太好看了","like":6789,"hate":29,"rcount":187,"floor":36,"ctime":"2024-12-15 11:05:00","mtime":"2024-12-15 11:05:00"}
{"index":{"_index":"replyrecord_00","_id":"30037_1037"}}
{"id":30037,"oid":1037,"mid":100,"type":1,"state":0,"content":"追更中","like":7890,"hate":34,"rcount":218,"floor":37,"ctime":"2024-12-15 11:06:00","mtime":"2024-12-15 11:06:00"}
{"index":{"_index":"replyrecord_00","_id":"30038_1038"}}
{"id":30038,"oid":1038,"mid":100,"type":1,"state":0,"content":"这就是天才吗","like":8901,"hate":39,"rcount":249,"floor":38,"ctime":"2024-12-15 11:07:00","mtime":"2024-12-15 11:07:00"}
{"index":{"_index":"replyrecord_00","_id":"30039_1039"}}
{"id":30039,"oid":1039,"mid":100,"type":1,"state":0,"content":"太强了太强了","like":9012,"hate":41,"rcount":267,"floor":39,"ctime":"2024-12-15 11:08:00","mtime":"2024-12-15 11:08:00"}
{"index":{"_index":"replyrecord_00","_id":"30040_1040"}}
{"id":30040,"oid":1040,"mid":100,"type":1,"state":0,"content":"完结撒花，期待下一部","like":10123,"hate":46,"rcount":289,"floor":40,"ctime":"2024-12-15 11:09:00","mtime":"2024-12-15 11:09:00"}
BULK_END

curl -s -X POST "$ES_HOST/_bulk" -H "Content-Type: application/x-ndjson" --data-binary @/tmp/replyrecord_bulk.json > /dev/null
echo "  评论记录数据写入完成 (40条)"

# 刷新索引
echo ""
echo "5. 刷新索引..."
curl -s -X POST "$ES_HOST/dm_search_000/_refresh" > /dev/null
curl -s -X POST "$ES_HOST/dm_date_2024_12/_refresh" > /dev/null
curl -s -X POST "$ES_HOST/pgc_media/_refresh" > /dev/null
curl -s -X POST "$ES_HOST/replyrecord_00/_refresh" > /dev/null
echo "  索引刷新完成"

# 验证数据
echo ""
echo "6. 验证数据..."
echo ""

dm_count=$(curl -s "$ES_HOST/dm_search_000/_count" | grep -o '"count":[0-9]*' | cut -d: -f2)
echo "  dm_search_000 文档数: $dm_count"

dm_date_count=$(curl -s "$ES_HOST/dm_date_2024_12/_count" | grep -o '"count":[0-9]*' | cut -d: -f2)
echo "  dm_date_2024_12 文档数: $dm_date_count"

pgc_count=$(curl -s "$ES_HOST/pgc_media/_count" | grep -o '"count":[0-9]*' | cut -d: -f2)
echo "  pgc_media 文档数: $pgc_count"

reply_count=$(curl -s "$ES_HOST/replyrecord_00/_count" | grep -o '"count":[0-9]*' | cut -d: -f2)
echo "  replyrecord_00 文档数: $reply_count"

# 清理临时文件
rm -f /tmp/dm_search_bulk.json /tmp/dm_date_bulk.json /tmp/pgc_media_bulk.json /tmp/replyrecord_bulk.json

echo ""
echo "=========================================="
echo "所有测试数据写入完成！"
echo "=========================================="
echo ""
echo "测试搜索示例:"
echo ""
echo "1. 弹幕搜索 (搜索包含'高能'的弹幕):"
echo "   curl -s '$ES_HOST/dm_search_000/_search' -H 'Content-Type: application/json' -d '{\"query\":{\"match\":{\"content\":\"高能\"}}}' | jq '.hits.hits[]._source.content'"
echo ""
echo "2. 弹幕日期搜索 (搜索 oid=1000 的弹幕日期统计):"
echo "   curl -s '$ES_HOST/dm_date_2024_12/_search' -H 'Content-Type: application/json' -d '{\"query\":{\"term\":{\"oid\":1000}}}' | jq '.hits.hits[]._source'"
echo ""
echo "3. PGC 番剧搜索 (搜索标题包含'进击'的番剧):"
echo "   curl -s '$ES_HOST/pgc_media/_search' -H 'Content-Type: application/json' -d '{\"query\":{\"match\":{\"title\":\"进击\"}}}' | jq '.hits.hits[]._source.title'"
echo ""
echo "4. 评论记录搜索 (搜索 mid=100 的评论):"
echo "   curl -s '$ES_HOST/replyrecord_00/_search' -H 'Content-Type: application/json' -d '{\"query\":{\"term\":{\"mid\":100}}}' | jq '.hits.total.value'"
echo ""
