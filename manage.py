#!/usr/bin/env python

import click

from app import models
from app.models import db
from app.main import app, setup_app


def abort_cb(ctx, param, value):
    if not value:
        ctx.abort()


def configure_cb(ctx, param, value):
    if value:
        click.secho("WARNING: Running in production mode", fg='yellow')
        app.config.from_object('config.Production')
    else:
        app.config.from_object('config.Development')

    setup_app()


@click.group()
@click.option('--production',
              is_flag=True,
              expose_value=False,
              is_eager=True,
              callback=configure_cb,
              help='Run with production configuration.')
def cli():
    pass


@cli.command()
def initdb():
    """Create tables for all models."""
    with app.app_context():
        db.create_all()
    click.echo('Initialized the database')


@cli.command()
@click.option('--yes', is_flag=True, callback=abort_cb,
              expose_value=False, prompt='Do you really want to drop the database?')
def dropdb():
    """Drop the database"""
    with app.app_context():
        db.drop_all()
    click.echo('Dropped the database')


@cli.command()
@click.option('--ipython/--no-ipython', default=True,
              help='Use IPython if it is installed')
def shell(ipython):
    """Runs a Python shell inside Flask application context."""
    context = dict(app=app, db=db, models=models)
    banner = '** Flask context shell **\n'

    if ipython:
        # Try IPython
        try:
            try:
                # 0.10.x
                from IPython.Shell import IPShellEmbed
                ipshell = IPShellEmbed(banner=banner)
                ipshell(global_ns=dict(), local_ns=context)

            except ImportError:
                # 0.12+
                from IPython import embed
                embed(banner1=banner, user_ns=context)
                return

        except ImportError:
            pass

    # Fall back to the built-in Python shell
    import code
    code.interact(banner, local=context)


@cli.command()
@click.option('--port', default=8080, help='Port to run the application on')
def run(port):
    """Run the application."""
    app.run(port=port)


if __name__ == "__main__":
    cli()
