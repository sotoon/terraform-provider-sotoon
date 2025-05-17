// Request Create Subnet Body No Internet
{"apiVersion":"networking.cafebazaar.cloud/v1alpha1","kind":"Subnet","metadata":{"name":"test"},"spec":{"cidr":"172.21.0.0/24"}}
// Response
{
  "apiVersion": "networking.cafebazaar.cloud/v1alpha1",
  "kind": "Subnet",
  "metadata": {
    "creationTimestamp": "2024-11-15T21:51:16Z",
    "generation": 1,
    "managedFields": [
      {
        "apiVersion": "networking.cafebazaar.cloud/v1alpha1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:spec": {
            ".": {},
            "f:cidr": {}
          }
        },
        "manager": "Mozilla",
        "operation": "Update",
        "time": "2024-11-15T21:51:16Z"
      }
    ],
    "name": "test",
    "namespace": "cafebazaar",
    "resourceVersion": "3198039656",
    "uid": "fbaf90a1-2f03-4d8a-b370-af18885d6247"
  },
  "spec": {
    "cidr": "172.21.0.0/24"
  }
}

// Request Create Subnet Body With Internet
{"apiVersion":"networking.cafebazaar.cloud/v1alpha1","kind":"Subnet","metadata":{"name":"test"},"spec":{"cidr":"172.21.0.0/24","routes":[{"to":"0.0.0.0/0","via":{"externalIP":"test"}}]}}

// Response
{
  "apiVersion": "networking.cafebazaar.cloud/v1alpha1",
  "kind": "Subnet",
  "metadata": {
    "creationTimestamp": "2024-11-15T21:52:52Z",
    "generation": 1,
    "managedFields": [
      {
        "apiVersion": "networking.cafebazaar.cloud/v1alpha1",
        "fieldsType": "FieldsV1",
        "fieldsV1": {
          "f:spec": {
            ".": {},
            "f:cidr": {},
            "f:routes": {}
          }
        },
        "manager": "Mozilla",
        "operation": "Update",
        "time": "2024-11-15T21:52:52Z"
      }
    ],
    "name": "test",
    "namespace": "cafebazaar",
    "resourceVersion": "3198041808",
    "uid": "d9e54c17-9532-4d14-a91e-2ea3279d7a1e"
  },
  "spec": {
    "cidr": "172.21.0.0/24",
    "routes": [
      {
        "to": "0.0.0.0/0",
        "via": {
          "externalIP": "test"
        }
      }
    ]
  }
}

// Request Create VM
{
    "apiVersion": "compute.ravh.ir/v1",
    "kind": "InstanceClaim",
    "metadata": {
        "name": "testvm"
    },
    "spec": {
        "iamEnabled": true,
        "subnetName": "default",
        "username": "compute",
        "type": "b2",
        "name": "testvm",
        "image": "ubuntu-2004-server",
        "disks": [
            {
                "remoteDisk": {
                    "gbSize": 20,
                    "name": "testvm-osdisk",
                    "tier": "ultra"
                }
            }
        ],
        "linkExternalIP": {
            "name": "testvm-ip0"
        }
    }
}

// Request Create IP
  {
    "apiVersion": "networking.cafebazaar.cloud/v1alpha1",
    "kind": "ExternalIP",
    "metadata": {
      "name": "testtest"
    },
    "spec": {
      "reserved": false
    }
  }

// Response
{
    "apiVersion": "networking.cafebazaar.cloud/v1alpha1",
    "kind": "ExternalIP",
    "metadata": {
      "creationTimestamp": "2024-11-16T08:02:22Z",
      "generation": 1,
      "managedFields": [
        {
          "apiVersion": "networking.cafebazaar.cloud/v1alpha1",
          "fieldsType": "FieldsV1",
          "fieldsV1": {
            "f:spec": {
              ".": {},
              "f:reserved": {}
            }
          },
          "manager": "Mozilla",
          "operation": "Update",
          "time": "2024-11-16T08:02:22Z"
        }
      ],
      "name": "testtest",
      "namespace": "cafebazaar",
      "resourceVersion": "3198787296",
      "uid": "e7999a19-e6cc-4b24-94ea-decda335f8f7"
    },
    "spec": {
      "reserved": false
    }
  }

// Patch IP
[{"op":"replace","path":"/spec/reserved","value":false}]

// Response
{
  "apiVersion": "networking.cafebazaar.cloud/v1alpha1",
  "kind": "ExternalIP",
  "metadata": {
    "annotations": {
      "release-version": "v0.1.0"
    },
    "creationTimestamp": "2022-07-13T14:47:37Z",
    
    "generation": 11,
    "name": "proxy1-eth0",
    "namespace": "cafebazaar",
    "resourceVersion": "3198818028",
    "uid": "3b0a3107-4508-4f35-865a-9ce6acf7b5e2"
  },
  "spec": {
    "gateway": "87.247.184.1",
    "ip": "87.247.186.8",
    "nodeName": "node0123.compute.neda",
    "poolName": "default",
    "rateLimiterName": "externalips-default",
    "reserved": false,
    "siteName": "default"
  },
  "status": {
    "ready": true
  }
}