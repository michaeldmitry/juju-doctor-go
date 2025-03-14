def init():
    jujudoctor.observe('status', on_status)

def on_status(event):
    """Handler for juju status event."""
    status = event.input
    agents = {}
    for app_name, app_info in status["applications"].items():
        if app_info["charm"] == "grafana-agent":
            agents[app_name] = app_info["subordinate-to"]


    one_grafana_agent_per_machine(event.input, agents)
    one_grafana_agent_per_app(event.input, agents)

def one_grafana_agent_per_machine(status, agents):
    # A mapping from grafana-agent app name to the list of apps it's subordinate to
    
    for agent, principals in agents.items():
        # A mapping from app name to machines
        machines = {}
        for p in principals:
            units = status["applications"][p].get("units", {})
            machines[p] = [unit["machine"] for unit in units.values()]

        for i in range(len(principals)):
            for j in range(i + 1, len(principals)):
                p1 = principals[i]
                p2 = principals[j]
                overlap = set(machines.get(p1, [])) & set(machines.get(p2, []))

                if overlap:
                    fail("{} is subordinate to both '{}' and '{}' in the same machines {}".format(
                        agent, p1, p2, overlap
                    ))



def one_grafana_agent_per_app(status, agents):
    # A mapping from grafana-agent app name to the list of apps it's subordinate to

    for agent, principals in agents.items():
        for p in principals:
            app_info = status["applications"].get(p, {})
            units = app_info.get("units", {})

            for name, unit in units.items():
                subord_apps = {}
                for sub in unit.get("subordinates", {}).keys():
                    subord_apps[sub.split("/", -1)[0]] = True  # Use dict as a set substitute

                # Find overlapping Grafana Agents
                subord_agents = []
                for sub_app in subord_apps:
                    if sub_app in agents:
                        subord_agents.append(sub_app)

                if len(subord_agents) > 1:
                    fail("{} is related to more than one grafana-agent subordinate: {}".format(
                        name, subord_agents
                    ))
