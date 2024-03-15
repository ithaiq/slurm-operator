# Slurm Operator使用指南

## 功能支持
- [x] jupyter(连通slurm集群下发任务)
- [x] slurmctld(node节点支持多实例)
- [x] slurmd(master单节点)
- [x] 统一挂载PVC
- [x] 支持自定义image、label、nodeSelector、resources等资源

## 编译
```bash
#指定版本打包编译docker镜像并推送
make docker-build docker-push IMG=harbor.xxx.cn:8443/operator/slurm-operator:v1.0.0
```

## 部署
```bash
#修改deploy/helm/slurm-operator/values.yaml参数执行以下命令(默认安装在kube-system命名空间)
helm install slurm-operator ./slurm-operator
#测试
kubectl apply -f config/samples/slurmoperator_v1beta1_slurmapplication.yaml
```

## 卸载
```bash
helm uninstall slurm-operator
```

## jupyter测试代码 (测试4个slurmd节点)
### python (test.py)
```python
#!/usr/bin/env python3
  
import time
import os
import socket
from datetime import datetime as dt
if __name__ == '__main__':
    print('Process started {}'.format(dt.now()))
    print('NODE : {}'.format(socket.gethostname()))
    print('PID  : {}'.format(os.getpid()))
    print('Executing for 15 secs')
    time.sleep(15)
    print('Process finished {}\n'.format(dt.now()))
```
### bash (job.sh)
```bash
#!/bin/bash
#
#SBATCH --job-name=test
#SBATCH --output=result.out
#
#SBATCH --ntasks=4
#SBATCH --ntasks-per-node=2
#
sbcast -f test.py /tmp/test.py
srun --nodes=4 python3 /tmp/test.py
```