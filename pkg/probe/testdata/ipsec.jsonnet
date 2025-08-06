# api/v2/monitor/vpn/ipsec?vdom=*
[
  {
    "http_method":"GET",
    "results":[
      {
        "proxyid":[
          {
            "proxy_src":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "proxy_dst":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "status":"up",
            "p2name":"tunnel_1-sub",
            "p2serial":1,
            "expire":11279,
            "incoming_bytes": 14298240,
            "outgoing_bytes": 14248560
          },
          {
            "proxy_src":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "proxy_dst":[
              {
                "subnet":"0.0.0.0\/0.0.0.0",
                "port":0,
                "protocol":0,
                "protocol_name":""
              }
            ],
            "status":"down",
            "p2name":"tunnel_1-sub",
            "p2serial":12,
            "expire":11279,
            "incoming_bytes": 14298240,
            "outgoing_bytes": 14248560
          }                    
        ],
        "name":"tunnel_1",
        "comments":"",
        "wizard-type":"custom",
        "creation_time":270801,
        "type":"automatic",
        "incoming_bytes": 14298240,
        "outgoing_bytes": 14248560,
        "rgwy":"1.2.3.4"
      },
      {
				"proxyid": [
					{
						"proxy_src": [
							{
								"subnet": "10.2.1.1",
								"port": 0,
								"protocol": 0,
								"protocol_name": ""
							}
						],
						"proxy_dst": [
							{
								"subnet": "10.2.1.1",
								"port": 0,
								"protocol": 0,
								"protocol_name": ""
							}
						],
						"status": "up",
						"p2name": "H-DC00_ISP1",
						"p2serial": 2,
						"expire": 3445,
						"incoming_bytes": 18961,
						"outgoing_bytes": 37172
					},
					{
						"proxy_src": [
							{
								"subnet": "0.0.0.0/0.0.0.0",
								"port": 0,
								"protocol": 0,
								"protocol_name": ""
							}
						],
						"proxy_dst": [
							{
								"subnet": "0.0.0.0/0.0.0.0",
								"port": 0,
								"protocol": 0,
								"protocol_name": ""
							}
						],
						"status": "up",
						"p2name": "H-DC00_ISP1",
						"p2serial": 1,
						"expire": 3438,
						"incoming_bytes": 13895,
						"outgoing_bytes": 38380
					}
				],
				"name": "H-DC00_ISP1_f",
				"parent": "H-DC00_ISP1",
				"comments": "",
				"wizard-type": "custom",
				"connection_count": 26,
				"creation_time": 717570,
				"username": "fwname_ISP1",
				"type": "dialup",
				"incoming_bytes": 68815634,
				"outgoing_bytes": 68824965,
				"rgwy": "10.1.1.1",
				"tun_id": "10.1.1.1",
				"tun_id6": "::10.1.1.1",
				"dialup_index": 15
			}
    ],
    "vdom":"root",
    "path":"vpn",
    "name":"ipsec",
    "status":"success",
    "serial":"FGT61FT000000000",
    "version":"v6.0.10",
    "build":365
  },
]
