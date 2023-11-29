import json

def make_targets():
    global uids
    with open('_task/wb_l0_data.json', encoding="utf8") as j:
        data = json.load(j)
        uids = [o['order_uid'] for o in data]
    
    with open('bench/target_cache.list', 'w+') as t:
        t.writelines(
            [f"GET http://localhost/v1/orders/{uid}\n" for uid in uids])

make_targets()