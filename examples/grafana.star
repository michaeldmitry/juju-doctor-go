def init():
    jujudoctor.observe('status', on_status)
    jujudoctor.observe('bundle', on_bundle)
    jujudoctor.observe('show_unit', on_show_unit)

def on_status(event):
    print("Ran STATUS probe: {}".format(event.input.keys()))

def on_bundle(event):
    print("Ran BUNDLE probe: {}".format(event.input.keys()))

def on_show_unit(event):
    print("Ran SHOW-UNIT probe: {}".format(event.input.keys()))